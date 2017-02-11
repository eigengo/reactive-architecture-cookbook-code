#include <algorithm>
#include <sstream>
#include <iostream>
#include <iomanip>
#include <random>
#include <experimental/optional>
#include <boost/asio.hpp>
#include <boost/program_options.hpp>
#include "faceextractor.h"
#include <librdkafka/rdkafkacpp.h>
#include <envelope.pb.h>
#include <faceextract-v1m0.pb.h>
#include <boost/uuid/uuid_generators.hpp>
#include <boost/uuid/uuid_io.hpp>
#include <google/protobuf/util/json_util.h>
#include <boost/algorithm/string/join.hpp>
#include <easylogging++.h>

#define DEMO

using namespace com::reactivearchitecturecookbook;

namespace asio = boost::asio;
namespace po = boost::program_options;
namespace fe = faceextract::v1m0;

template<typename T, typename U, typename Handler>
std::vector<Envelope> handle_sync(const std::unique_ptr<RdKafka::Message> &message, Handler handler) {
    Envelope envelope;
    if (envelope.ParseFromArray(message->payload(), static_cast<int>(message->len()))) {
        if (envelope.payload().Is<T>()) {
            T payload;
            if (envelope.payload().UnpackTo(&payload)) {
                std::vector<U> mapped_payloads = handler(envelope, payload);
                std::vector<Envelope> result(mapped_payloads.size());
                std::transform(mapped_payloads.begin(), mapped_payloads.end(), result.begin(),
                               [&envelope](const U &mapped_payload) {
                                   Envelope mapped_envelope(envelope);
                                   mapped_envelope.mutable_payload()->PackFrom(mapped_payload);
                                   return mapped_envelope;
                               });
                return result;
            }
        }
    }
    return std::vector<Envelope>();
}

static std::atomic_bool run(true);
static void sigterm(int) {
    run = false;
}

class commit_dr_cb : public RdKafka::DeliveryReportCb {
private:
    std::shared_ptr<RdKafka::KafkaConsumer> consumer_;
public:
    commit_dr_cb(std::shared_ptr<RdKafka::KafkaConsumer> consumer) : consumer_(consumer) { };

    void dr_cb(RdKafka::Message &message) override {
        if (message.err() != RdKafka::ERR_NO_ERROR) return;

        auto *tp = static_cast<RdKafka::TopicPartition *>(message.msg_opaque());
        if (tp) {
            LOG(INFO) << "Committed offset " << tp->offset() << " for " << message.key() << " on " << message.topic_name();
            consumer_->commitAsync(std::vector<RdKafka::TopicPartition *>{tp});
            delete tp;
        }
    }
};

INITIALIZE_EASYLOGGINGPP

int main(int argc, const char *argv[]) {
    std::vector<std::string> brokers{"localhost"};
    std::string in_topic_name = "ingest-1";
    std::string out_topic_name = "vision-1";
    std::string consumer_group = "recognition";

    po::options_description description;
    description.add_options()
            ("brokers,B", po::value<std::vector<std::string>>(&brokers), "Set the brokers bindings")
            ("in-topic,I", po::value<std::string>(&in_topic_name), "Set the input topic")
            ("out-topic,O", po::value<std::string>(&out_topic_name), "Set the output topic")
            ("group,G", po::value<std::string>(&consumer_group), "Set the consumer group");

    po::command_line_parser(argc, argv).options(description).run();

    po::variables_map vm;
    po::store(po::command_line_parser(argc, argv).options(description).run(), vm);
    po::notify(vm);

    const int32_t partition = RdKafka::Topic::PARTITION_UA;
    auto conf = std::unique_ptr<RdKafka::Conf>(RdKafka::Conf::create(RdKafka::Conf::CONF_GLOBAL));
    auto tconf = std::unique_ptr<RdKafka::Conf>(RdKafka::Conf::create(RdKafka::Conf::CONF_TOPIC));
    std::string err_str;
    conf->set("metadata.broker.list", boost::algorithm::join(brokers, ","), err_str);
    conf->set("group.id", consumer_group, err_str);

    auto consumer = std::shared_ptr<RdKafka::KafkaConsumer>(RdKafka::KafkaConsumer::create(conf.get(), err_str));
    commit_dr_cb dr_cb(consumer);
    conf->set("dr_cb", &dr_cb, err_str);
    auto producer = std::unique_ptr<RdKafka::Producer>(RdKafka::Producer::create(conf.get(), err_str));
    const auto out_topic = std::unique_ptr<RdKafka::Topic>(
            RdKafka::Topic::create(producer.get(), out_topic_name, tconf.get(), err_str));
    consumer->subscribe(std::vector<std::string>({in_topic_name}));

    signal(SIGINT, sigterm);
    signal(SIGTERM, sigterm);

    while (run) {
        auto message = std::unique_ptr<RdKafka::Message>(consumer->consume(1000));
#ifdef DEMO
        boost::uuids::random_generator uuid_gen;
        std::random_device rd;
        std::mt19937 key_gen(rd());
        std::uniform_int_distribution<> dis(0, 5);
        const std::string k = std::to_string(dis(key_gen));
        const auto key = &k;
        Envelope envelope;
        envelope.set_correlation_id(boost::uuids::to_string(uuid_gen()));
        fe::ExtractedFace ef;
        envelope.mutable_payload()->PackFrom(ef);
        const auto out_envelopes = std::vector<Envelope>{ envelope };
#else
        if (message->err() != RdKafka::ERR_NO_ERROR) continue;
        const auto &key = message->key();
        const auto out_envelopes = handle_sync<fe::ExtractFace, fe::ExtractedFace>(
                message, [](const auto &, const auto &) { return std::vector<fe::ExtractedFace>{fe::ExtractedFace()}; });
#endif
        for (size_t i = 0; i < out_envelopes.size(); ++i) {
            const auto &out_envelope = out_envelopes[i];
            RdKafka::TopicPartition *opaque = nullptr;
            if (i == out_envelopes.size() - 1) {
                opaque = RdKafka::TopicPartition::create(message->topic_name(), message->partition(), message->offset());
            }
            const auto out_payload = out_envelope.SerializeAsString();
            const auto resp = producer->produce(out_topic.get(), partition,
                                                RdKafka::Producer::RK_MSG_COPY,
                                                const_cast<char *>(out_payload.c_str()), out_payload.size(),
                                                key, opaque);
            if (resp != RdKafka::ERR_NO_ERROR) {
                delete opaque;
                LOG(ERROR) << "Produce failed: " << RdKafka::err2str(resp);
            } else {
                LOG(INFO) << "Produced message (" << out_payload.size() << " bytes) from " << message->offset();
            }
        }

        producer->poll(0);
    }

    consumer->close();
    producer->flush(1000);
}
