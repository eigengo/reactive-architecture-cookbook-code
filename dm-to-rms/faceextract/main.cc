#include <algorithm>
#include <sstream>
#include <iostream>
#include <iomanip>
#include <experimental/optional>
#include <boost/asio.hpp>
#include <boost/program_options.hpp>
#include "recogniser.h"
#include <librdkafka/rdkafkacpp.h>
#include <envelope.pb.h>
#include <faceextract-v1m0.pb.h>
#include <boost/uuid/uuid_generators.hpp>
#include <boost/uuid/uuid_io.hpp>
#include <google/protobuf/util/json_util.h>
#include <boost/algorithm/string/join.hpp>

using namespace com::reactivearchitecturecookbook;

namespace asio = boost::asio;
namespace po = boost::program_options;

template<typename T, typename U, typename Handler>
std::vector<Envelope> transform(RdKafka::Message *message, Handler mapper) {
    Envelope envelope;
    if (envelope.ParseFromArray(message->payload(), static_cast<int>(message->len()))) {
        if (envelope.payload().Is<T>()) {
            T payload;
            if (envelope.payload().UnpackTo(&payload)) {
                std::vector<U> mapped_payloads = mapper(payload);
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

    auto producer = std::unique_ptr<RdKafka::Producer>(RdKafka::Producer::create(conf.get(), err_str));
    const auto out_topic = std::unique_ptr<RdKafka::Topic>(
            RdKafka::Topic::create(producer.get(), out_topic_name, tconf.get(), err_str));
    auto consumer = std::unique_ptr<RdKafka::KafkaConsumer>(RdKafka::KafkaConsumer::create(conf.get(), err_str));
    const std::string key = "key";
    consumer->subscribe(std::vector<std::string>({in_topic_name}));

    signal(SIGINT, sigterm);
    signal(SIGTERM, sigterm);

    while (run) {
        const auto message = std::unique_ptr<RdKafka::Message>(consumer->consume(10));
        if (message->err() != RdKafka::ERR_NO_ERROR) continue;

        const auto out_envelopes = transform<faceextract::v1m0::ExtractFace, faceextract::v1m0::ExtractedFace>(
                message.get(),
                [](const auto &extractFace) {
                    return std::vector<faceextract::v1m0::ExtractedFace>{
                            faceextract::v1m0::ExtractedFace()};
                });
        for (const auto &out_envelope : out_envelopes) {
            const auto out_payload = out_envelope.SerializeAsString();
            const auto resp = producer->produce(out_topic.get(), partition,
                                                RdKafka::Producer::RK_MSG_COPY,
                                                const_cast<char *>(out_payload.c_str()), out_payload.size(),
                                                &key, nullptr);
            if (resp != RdKafka::ERR_NO_ERROR) {
                std::cerr << "% Produce failed: " << RdKafka::err2str(resp) << std::endl;
            } else {
                consumer->commitSync(message.get());
                std::cerr << "% Produced message (" << out_payload.size() << " bytes)" << std::endl;
            }
        }

        producer->poll(0);
    }

    consumer->close();
    producer->flush(0);
}
