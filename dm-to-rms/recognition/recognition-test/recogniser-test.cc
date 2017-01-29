#include <gtest/gtest.h>
#include <opencv/cv.hpp>
#include "recogniser.h"

using namespace com::reactivearchitecturecookbook;

class recognition_test : public testing::Test {
protected:
    recogniser recogniser;
};

TEST_F(recognition_test, recognise_trivial) {
    recogniser.recognise(cv::Mat());
}
