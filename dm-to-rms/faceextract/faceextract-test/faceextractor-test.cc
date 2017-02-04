#include <gtest/gtest.h>
#include <opencv/cv.hpp>
#include "faceextractor.h"

using namespace com::reactivearchitecturecookbook;

class faceextractor_test : public testing::Test {
protected:
    face_exctactor fe;
};

TEST_F(faceextractor_test, recognise_trivial) {
    fe.extract(cv::Mat());
}
