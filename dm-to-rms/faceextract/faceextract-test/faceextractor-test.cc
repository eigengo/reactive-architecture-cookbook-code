#include <gtest/gtest.h>
#include <opencv/cv.hpp>
#include "recogniser.h"

using namespace com::reactivearchitecturecookbook;

class faceextractor_test : public testing::Test {
protected:
    faceextractor faceextractor;
};

TEST_F(faceextractor_test, recognise_trivial) {
    faceextractor.extract(cv::Mat());
}
