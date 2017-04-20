#include <gtest/gtest.h>
#include <rapidcheck.h>
#include <rapidcheck/gtest.h>
#include "protobuf_gen.h"
#include <ingest-v1m0.pb.h>

using namespace com::reactivearchitecturecookbook;

RC_GTEST_PROP(main_test, handle_extract_face, (const ingest::v1m0::IngestedImage &gen)) {
    ingest::v1m0::IngestedImage ser;
    ser.ParseFromString(gen.SerializeAsString());

    RC_ASSERT(ser.content() == gen.content());
    RC_ASSERT(ser.media_type() == gen.media_type());

    std::cout << "mime_type = " << gen.media_type() << ", content = " << gen.content() << std::endl;
}
