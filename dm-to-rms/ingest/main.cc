#include <algorithm>
#include <sstream>
#include <iostream>
#include <iomanip>
#include <azmq/socket.hpp>
#include <boost/asio.hpp>
#include <boost/program_options.hpp>
#include <boost/uuid/uuid_generators.hpp>
#include <boost/uuid/uuid_io.hpp>
#include <google/protobuf/util/json_util.h>
#include <nghttp2/asio_http2_server.h>

namespace asio = boost::asio;
namespace po = boost::program_options;
namespace ah = nghttp2::asio_http2;
namespace bs = boost::system;

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wmissing-noreturn"

/**
 * Implements the main loop in the worker
 *
 * \param in the pulling socket
 * \param out the output pushing socket
 * \remark this function diverges
 */
[[noreturn]]
void main_loop() {
    while (true);
}

#pragma clang diagnostic pop

int main(int argc, const char *argv[]) {
    ah::server::http2 server;

    std::string style_css = "h1 { color: green; }";
    server.handle("/", [&style_css](const ah::server::request &request, const ah::server::response &response) {
        boost::system::error_code ec;
        auto push = response.push(ec, "GET", "/style.css");
        push->write_head(200);
        push->end(style_css);

        response.write_head(200);
        response.end(R"(
<!DOCTYPE html><html lang="en">
<title>HTTP/2 FTW</title><body>
<link href="/style.css" rel="stylesheet" type="text/css">
<h1>This should be green</h1>
</body></html>
)");
    });

    server.handle("/style.css", [&style_css](const ah::server::request &req, const ah::server::response &res) {
        res.write_head(200);
        res.end(style_css);
    });

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
