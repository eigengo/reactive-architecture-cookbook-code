#include <algorithm>
#include <sstream>
#include <iostream>
#include <iomanip>
#include <azmq/socket.hpp>
#include <boost/asio.hpp>
#include "recogniser.h"

#include <envelope.pb.h>
#include <faceextract-v1m0.pb.h>
#include <boost/uuid/uuid_generators.hpp>
#include <boost/uuid/uuid_io.hpp>
#include <google/protobuf/util/json_util.h>

using namespace com::reactivearchitecturecookbook;

/**
 * Convenience function that wraps the `payload` into a newly constructed envelope.
 * The newly created envelope starts with a randomly-generated correlation id.
 *
 * @param payload the payload to wrap
 * @return an envelope
 */
Envelope new_envelope_with_payload(google::protobuf::Message &payload) {
    Envelope envelope;
    auto uuid = boost::uuids::basic_random_generator<boost::mt19937>()();
    envelope.set_correlation_id(to_string(uuid));
    envelope.mutable_payload()->PackFrom(payload);
    return envelope;
}

/**
 * This is a copy-and-paste from SO: formats characters in `input` as their hex values
 * @param input the input string
 * @return the hex representation of `input`
 */
std::string string_to_hex(const std::string &input) {
    static const char *const lut = "0123456789ABCDEF";
    size_t len = input.length();

    std::string output;
    output.reserve(2 * len);
    for (size_t i = 0; i < len; ++i) {
        const char c = input[i];
        output.push_back(lut[c >> 4]);
        output.push_back(lut[c & 15]);
    }
    return output;
}

namespace asio = boost::asio;

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wmissing-noreturn"

[[noreturn]]
int main(int argc, char **argv) {
    asio::io_service ios1;
    azmq::sub_socket subscriber(ios1, ZMQ_DEALER);
    azmq::pub_socket publisher(ios1, ZMQ_ROUTER);
    publisher.bind("ipc://foo");
    subscriber.connect("ipc://foo");

    while (true) {
        std::cout << "." << std::endl;
        sleep(1);
        std::array<unsigned char, 1024> buf;
        publisher.send(asio::buffer(buf));
        subscriber.async_receive([](const boost::system::error_code, azmq::message, size_t) {
            std::cout << "x" << std::endl;
        });

        ios1.run();
    }
/*
    recogniser recogniser;
    cv::Mat image;
    recogniser.recognise(std::forward<cv::Mat>(image));
    std::cout << "r" << std::endl;

    // Use the generated class from the protocol definition
    faceextract::v1m0::ExtractFace extractFace;
    extractFace.set_mime_type("image/png");
    extractFace.set_content("<png-bytes>");

    // wrap it in an envelope
    auto msg = new_envelope_with_payload(extractFace);

    // serialize to json (rather than the binary, which is shown above)
    std::string json;
    google::protobuf::util::MessageToJsonString(msg, &json);

    // output the json
    std::cout << "json: " << json << std::endl;

    // output to raw bytes
    std::string raw;
    msg.SerializeToString(&raw);
    std::cout << "bytes: " << string_to_hex(raw) << std::endl;
*/
}

#pragma clang diagnostic pop
