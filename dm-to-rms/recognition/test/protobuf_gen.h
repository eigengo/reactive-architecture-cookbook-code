#ifndef ALL_PROTOBUF_GEN_H
#define ALL_PROTOBUF_GEN_H

#include <rapidcheck.h>
#include <faceextract-v1m0.pb.h>

using namespace com::reactivearchitecturecookbook;

// NOTE: Must be in rc namespace!
namespace rc {

    template<>
    struct Arbitrary<faceextract::v1m0::ExtractFace> {
        static Gen <faceextract::v1m0::ExtractFace> arbitrary() {
            const std::vector<std::string> contentTypes {"image/png", "image/jpeg", "image/tiff", "image/bmp"};
            const auto contentTypesGen2 = gen::map(gen::container<std::vector<char>>(32, gen::inRange<char>(32, 127)), [](const auto &chars) {
                return std::string(chars.begin(), chars.end());
            });
            const auto contentTypesGen = gen::elementOf(contentTypes);

            auto x = gen::tuple(contentTypesGen2, gen::arbitrary<std::string>());
            return gen::map(x, [](const auto pair) {
                faceextract::v1m0::ExtractFace extractFace;
                extractFace.set_mime_type(std::get<0>(pair));
                extractFace.set_content(std::get<1>(pair));
                return extractFace;
            });
        };
    };

};


#endif //ALL_PROTOBUF_GEN_H
