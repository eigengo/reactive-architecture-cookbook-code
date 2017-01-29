#include <gtest/gtest.h>
#include <rapidcheck.h>
#include <rapidcheck/gtest.h>
#include "protobuf_gen.h"
#include <envelope.pb.h>
#include <faceextract-v1m0.pb.h>

using namespace com::reactivearchitecturecookbook;

class main_test : public testing::Test {
protected:
};

RC_GTEST_FIXTURE_PROP(main_test, handle_extract_face, (const faceextract::v1m0::ExtractFace &gen)) {
    faceextract::v1m0::ExtractFace ser;
    ser.ParseFromString(gen.SerializeAsString());

    RC_ASSERT(ser.content() == gen.content());
    RC_ASSERT(ser.mime_type() == gen.mime_type());

    std::cout << "mime_type = " << gen.mime_type() << ", content = " << gen.content() << std::endl;
}


TEST_F(main_test, x) {
    faceextract::v1m0::ExtractFace extractFace;

    const google::protobuf::Descriptor *descriptor = extractFace.descriptor();
    for (int i = 0; i < descriptor->field_count(); ++i) {
        const google::protobuf::FieldDescriptor *fieldDescriptor = descriptor->field(i);
        switch (fieldDescriptor->cpp_type()) {
            case google::protobuf::FieldDescriptor::CppType::CPPTYPE_STRING:
                std::cout << fieldDescriptor->name() << std::endl;
                break;
            default: break;
        }
    }

}
