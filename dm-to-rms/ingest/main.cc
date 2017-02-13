#include <algorithm>
#include <sstream>
#include <iostream>
#include <iomanip>
#include <azmq/socket.hpp>
#include <boost/asio.hpp>
#include <boost/algorithm/string/join.hpp>
#include <boost/program_options.hpp>
#include <boost/uuid/uuid_generators.hpp>
#include <boost/uuid/uuid_io.hpp>
#include <google/protobuf/util/json_util.h>
#include <nghttp2/asio_http2_server.h>
#include <librdkafka/rdkafkacpp.h>
#include <envelope.pb.h>
#include <ingest-v1m0.pb.h>
#include <easylogging++.h>

using namespace com::reactivearchitecturecookbook;

namespace in = ingest::v1m0;
namespace asio = boost::asio;
namespace po = boost::program_options;
namespace ngs = nghttp2::asio_http2::server;
namespace bs = boost::system;

INITIALIZE_EASYLOGGINGPP

int main(int argc, const char *argv[]) {
    std::vector<std::string> brokers{"localhost"};
    std::string out_topic_name = "ingest-1";
    std::string host = "localhost";
    int port = 8000;

    po::options_description description;
    description.add_options()
            ("brokers,B", po::value<std::vector<std::string>>(&brokers), "Set the brokers bindings")
            ("out-topic,O", po::value<std::string>(&out_topic_name), "Set the output topic")
            ("host,H", po::value<std::string>(&host), "Set the HTTP/2 server host")
            ("port,P", po::value<int>(&port), "Set the HTTP/2 port");

    po::command_line_parser(argc, argv).options(description).run();

    po::variables_map vm;
    po::store(po::command_line_parser(argc, argv).options(description).run(), vm);
    po::notify(vm);

    ngs::http2 server;
    const int32_t partition = RdKafka::Topic::PARTITION_UA;
    auto conf = std::unique_ptr<RdKafka::Conf>(RdKafka::Conf::create(RdKafka::Conf::CONF_GLOBAL));
    auto tconf = std::unique_ptr<RdKafka::Conf>(RdKafka::Conf::create(RdKafka::Conf::CONF_TOPIC));
    std::string err_str;
    conf->set("metadata.broker.list", boost::algorithm::join(brokers, ","), err_str);

    auto producer = std::unique_ptr<RdKafka::Producer>(RdKafka::Producer::create(conf.get(), err_str));
    const auto out_topic = std::unique_ptr<RdKafka::Topic>(
            RdKafka::Topic::create(producer.get(), out_topic_name, tconf.get(), err_str));

    server.handle("/", [&](const ngs::request &request, const ngs::response &response) {
        auto authorisation = request.header().find("Authorization");
        request.on_data([&](const uint8_t *data, std::size_t size) {
            Envelope out_envelope;
            in::IngestedImage ingestedImage;
            ingestedImage.set_mime_type("image/png");

            out_envelope.set_token(
                    "eyJlbmMiOiJBMTI4R0NNIiwiYWxnIjoiUlNBLU9BRVAifQ.aJAcmDd3Do2CTI50zNvPRajU9NIwoiSgjKs9CkCeP_o9YpYxsemo5ij1WVdC3Cy4LEHian5m01kiiwpxV0PQHbCmobbFOnS1YuCs8LEkSLJLsFNZjw4DDRd2sB0UY9wXwEE0SVSx3u05VkV-w-haj5W9byaJ0OkNssFsZzzfzmc.Q0bDPb-LpTzNCqTw.Vo_dsNW1v-lxV7FgD-R9IrCMfTaQ-1ud07yhja0Ng2QtlmmoNVHtWjG_TP1WTIm-QuNnRFakvPwllIVNgtGDm9qecsImTJB2bOa-D1TEJOeNIKQLPChCdGm_cdHKLfn0zgLbH6q61EcaDT-GNmWvy2FkvdvRg_5P1XFueHCxK4s9Z1-sfxCR8T5J6Ej_WftaLrhlP0FMqDIfLg3_VGbodv-A8g3G_Taqnmd8eQt3JmVt9mn3MKBNfnejAtKaJwxXx-WGbwCOESz7gSNS6DuoccSsCrQqJ-VmnP_Kj9E785Fpum4Z3Vv0Htj6EABkgyPxoH3KlEy4zDRL_rIk7zTFXF67MUtPpb6NKcw.8_l5m0IeqxbqQbNeIPXgJg");
            out_envelope.mutable_payload()->PackFrom(ingestedImage);
            boost::uuids::random_generator uuid_gen;
            out_envelope.set_correlation_ids(0, boost::uuids::to_string(uuid_gen()));

            const std::string key = request.uri().path;
            const auto out_payload = out_envelope.SerializeAsString();
            const auto resp = producer->produce(out_topic.get(), partition,
                                                RdKafka::Producer::RK_MSG_COPY,
                                                const_cast<char *>(out_payload.c_str()), out_payload.size(),
                                                &key, nullptr);
            if (resp != RdKafka::ERR_NO_ERROR) {
                response.write_head(500);
                response.end("{'error':'" + RdKafka::err2str(resp) + "'}");
                LOG(ERROR) << "Produce failed: " << RdKafka::err2str(resp);
            } else {
                response.write_head(200);
                response.end("{'bytes':" + std::to_string(out_payload.size()) + "}");
                LOG(INFO) << "Produced message (" << out_payload.size() << " bytes)";
            }
        });
    });

    bs::error_code ec;
    asio::ssl::context tls(boost::asio::ssl::context::sslv23);

    tls.use_private_key_file("../server.key", boost::asio::ssl::context::pem);
    tls.use_certificate_chain_file("../server.crt");

    ngs::configure_tls_context_easy(ec, tls);
    if (server.listen_and_serve(ec, tls, "localhost", "8000")) {
        std::cerr << "error: " << ec.message() << std::endl;
    }
}
