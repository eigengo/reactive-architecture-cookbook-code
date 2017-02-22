#ifndef ALL_PROTOBUF_GEN_H
#define ALL_PROTOBUF_GEN_H

#include <rapidcheck.h>
#include <ingest-v1m0.pb.h>

using namespace com::reactivearchitecturecookbook;

#include <type_traits>
template<typename T>
class protobuf_gen {
    static_assert(std::is_base_of<google::protobuf::Message, T>::value, "T must extend google::protobuf::Message");
public:

    rc::Gen<T> pb() {
        return [](const rc::Random &random, int size) {
            T gen(T::default_instance());
            gen.GetReflection()->SetString(&gen, gen.GetDescriptor()->FindFieldByName("x"), "");

            return rc::gen::just(gen);
        };
    }


};


// NOTE: Must be in rc namespace!
namespace rc {

    template<>
    struct Arbitrary<ingest::v1m0::IngestedImage> {
        static Gen<ingest::v1m0::IngestedImage> arbitrary() {
            const std::vector<std::string> contentTypes {"image/png", "image/jpeg", "image/tiff", "image/bmp"};
            const auto contentTypesGen2 = gen::map(gen::container<std::vector<char>>(gen::inRange<char>(32, 127)), [](const auto &chars) {
                return std::string(chars.begin(), chars.end());
            });
            const auto contentTypesGen = gen::elementOf(contentTypes);

            auto x = gen::tuple(contentTypesGen2, gen::arbitrary<std::string>());

            return gen::map(x, [](const auto pair) {
                ingest::v1m0::IngestedImage ingested_image;
                ingested_image.set_mime_type(std::get<0>(pair));
                ingested_image.set_content(std::get<1>(pair));
                return ingested_image;
            });
        };
    };

};


#endif //ALL_PROTOBUF_GEN_H
