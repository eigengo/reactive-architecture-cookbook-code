#include <gtest/gtest.h>
#include <rapidcheck/gtest.h>
#include "protobuf_gen.h"
#include <envelope.pb.h>
#include <faceextract-v1m0.pb.h>

using namespace com::reactivearchitecturecookbook;

class main_test : public testing::Test {
};

TEST_F(main_test, x) {
    faceextract::v1m0::ExtractFace extractFace;
    const google::protobuf::Descriptor *descriptor = extractFace.descriptor();
    for (int i = 0; i < descriptor->field_count(); ++i) {
        const google::protobuf::FieldDescriptor *fieldDescriptor = descriptor->field(i);
        switch (fieldDescriptor->cpp_type()) {
            case google::protobuf::FieldDescriptor::CppType::CPPTYPE_STRING: break;
            default: break;
        }
    }
}

RC_GTEST_FIXTURE_PROP(main_test,
                      copyOfStringIsIdenticalToOriginal,
                      (const std::string &str)) {
    const auto strCopy = str;
    RC_ASSERT(strCopy == str);
}
