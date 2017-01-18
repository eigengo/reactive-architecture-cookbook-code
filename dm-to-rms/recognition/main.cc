#include <algorithm>
#include <sstream>
#include <iostream>
#include <iomanip>
#include <azmq/socket.hpp>
#include <boost/asio.hpp>
#include <boost/program_options.hpp>
#include "recogniser.h"

#include <envelope.pb.h>
#include <faceextract-v1m0.pb.h>
#include <boost/uuid/uuid_generators.hpp>
#include <boost/uuid/uuid_io.hpp>
#include <google/protobuf/util/json_util.h>

using namespace com::reactivearchitecturecookbook;


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
void main_loop(azmq::pull_socket &in, azmq::push_socket &out) {
    while (true) {
        azmq::message in_data;
        in.receive(in_data);
        std::cout << "." << std::endl;
        Envelope envelope;
        envelope.ParseFromArray(in_data.data(), static_cast<int>(in_data.size()));
        if (envelope.payload().Is<faceextract::v1m0::ExtractFace>()) {
            faceextract::v1m0::ExtractFace extractFace;
            envelope.payload().UnpackTo(&extractFace);

            std::cout << "Extracting from " << extractFace.mime_type() << std::endl;

            faceextract::v1m0::ExtractedFace extractedFace;

            auto response = Envelope(envelope);
            response.mutable_payload()->PackFrom(extractedFace);
            auto out_data = asio::buffer(response.SerializeAsString());
            out.send(out_data);
        } else {
            out.send(azmq::message("boo!"));
        }
    }
}

#pragma clang diagnostic pop

int main(int argc, const char *argv[]) {
    asio::io_service ios;
    azmq::pull_socket in(ios);
    azmq::push_socket out(ios);

    std::vector<std::string> in_addresses{"tcp://localhost:5555"};
    std::vector<std::string> out_bindings{"tcp://*:5556"};

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

    main_loop(in, out);
}

