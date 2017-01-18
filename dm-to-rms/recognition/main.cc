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

Envelope reply_to(Envelope &envelope, google::protobuf::Message &payload) {
    Envelope reply_envelope;
    reply_envelope.set_correlation_id(envelope.correlation_id());
    reply_envelope.mutable_payload()->PackFrom(payload);
    return reply_envelope;
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
/**
 * Implements the main loop in the worker
 */
[[noreturn]]
void main_loop(azmq::pull_socket &in, azmq::push_socket &out) {
    while (true) {
        azmq::message in_data;
        in.receive(in_data);
        std::cout << "." << std::endl;
        Envelope envelope;
        envelope.ParseFromArray(in_data.data(), in_data.size());
        if (envelope.payload().Is<faceextract::v1m0::ExtractFace>()) {
            faceextract::v1m0::ExtractFace extractFace;
            envelope.payload().UnpackTo(&extractFace);

            std::cout << "Extracting from " << extractFace.mime_type() << std::endl;

            faceextract::v1m0::ExtractedFace extractedFace;
            auto response = reply_to(envelope, extractedFace);
            auto out_data = asio::buffer(response.SerializeAsString());
            out.send(out_data);
        }
    }
}

#pragma clang diagnostic pop

int main(int argc, char **argv) {
    asio::io_service ios;
    azmq::pull_socket in(ios);
    azmq::push_socket out(ios);

    in.connect("tcp://localhost:5555");
    out.bind("tcp://*:5556");

    main_loop(in, out);
/*

    asio::io_service ios_s, ios_p;
    azmq::sub_socket subscriber(ios_s);
    azmq::pub_socket publisher(ios_p);
    publisher.bind("tcp:// *:5556");
    subscriber.connect("tcp://localhost:5556");
    subscriber.set_option(azmq::socket::subscribe(""));

    while (true) {
        sleep(1);
        auto sent = publisher.send(azmq::message("fpop"));
        std::cout << "sent " << sent << std::endl;
        azmq::message received;
        subscriber.receive(received);
        std::cout << received.string() << std::endl;
    }

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

