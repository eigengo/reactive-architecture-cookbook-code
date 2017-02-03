#include <algorithm>
#include <sstream>
#include <iostream>
#include <iomanip>
#include <boost/asio.hpp>
#include <boost/program_options.hpp>
#include <boost/uuid/uuid_generators.hpp>
#include <boost/uuid/uuid_io.hpp>
#include <boost/filesystem.hpp>
#include <google/protobuf/util/json_util.h>
#include <nghttp2/asio_http2_server.h>
#include <libkafka_asio/libkafka_asio.h>

namespace asio = boost::asio;
namespace po = boost::program_options;
namespace ah = nghttp2::asio_http2;
namespace bs = boost::system;
namespace fs = boost::filesystem;
namespace kafka = libkafka_asio;

int main(int argc, const char *argv[]) {
    ah::server::http2 server;
    kafka::Connection::Configuration configuration;
    std::string topic_name = "ingest-1.0.0";
    fs::path server_key, server_cert;
    configuration.auto_connect = true;
    configuration.client_id = "ingest-1.0.0";
    configuration.socket_timeout = 1000;
    kafka::ConnectionConfiguration::BrokerAddress broker_address{.hostname = "192.168.0.7", .service = "2181"};
    configuration.

    asio::io_service kafka_ios;
    kafka::Connection connection(kafka_ios, configuration);
    connection.AsyncConnect([](const boost::system::error_code& error) {
        std::cout << error << std::endl;
    });

    server.handle("/", [&connection, &topic_name](const auto &http_req, const auto &http_resp) {
        kafka::ProduceRequest kafka_req;
        kafka_req.AddValue("foo", topic_name);
        std::cout << "->K" << std::endl;
        connection.AsyncRequest(kafka_req, [&http_resp](const auto &err, const auto &kafka_resp) {
            std::cout << "=>K" << std::endl;
            if (err) http_resp.write_head(500);
            else http_resp.write_head(200);
            http_resp.end("done");
        });
    });

//    std::string style_css = "h1 { color: green; }";
//    server.handle("/index.html", [&style_css](const ah::server::request &request, const ah::server::response &response) {
//        boost::system::error_code ec;
//        auto push = response.push(ec, "GET", "/style.css");
//        push->write_head(200);
//        push->end(style_css);
//
//        response.write_head(200);
//        response.end(R"(
//<!DOCTYPE html><html lang="en">
//<title>HTTP/2 FTW</title><body>
//<link href="/style.css" rel="stylesheet" type="text/css">
//<h1>This should be green</h1>
//</body></html>
//)");
//    });
//
//    server.handle("/style.css", [&style_css](const ah::server::request &req, const ah::server::response &res) {
//        res.write_head(200);
//        res.end(style_css);
//    });

    bs::error_code ec;
    asio::ssl::context tls(boost::asio::ssl::context::sslv23);

    tls.use_private_key_file("../server.key", boost::asio::ssl::context::pem);
    tls.use_certificate_chain_file("../server.crt");

    ah::server::configure_tls_context_easy(ec, tls);
    if (server.listen_and_serve(ec, tls, "localhost", "8000")) {
        std::cerr << "error: " << ec.message() << std::endl;
    }

    /*
      std::vector<std::string> in_addresses{"tcp://localhost:5555"};
      std::vector<std::string> out_bindings{"tcp:// *:5556"};

      po::options_description description;
      description.add_options()
              ("in,I",  po::value<std::vector<std::string>>(&in_addresses), "Set the input bindings")
              ("out,O", po::value<std::vector<std::string>>(&out_bindings), "Set the output bindings");

      po::command_line_parser(argc, argv).options(description).run();

      po::variables_map vm;
      po::store(po::command_line_parser(argc, argv).options(description).run(), vm);
      po::notify(vm);

      std::for_each(in_addresses.begin(), in_addresses.end(), [&in](auto a) { in.connect(a); });
      std::for_each(out_bindings.begin(), out_bindings.end(), [&out](auto a) { out.bind(a); });
  */
    // main_loop();
}
