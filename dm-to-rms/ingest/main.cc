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

namespace asio = boost::asio;
namespace po = boost::program_options;

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
    while (true) ;
}

#pragma clang diagnostic pop

int main(int argc, const char *argv[]) {
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
    main_loop();
}
