#include <string>
#include <memory>
#include <vector>
#include <opencv/cv.hpp>

namespace com {
    namespace reactivearchitecturecookbook {

        class face_exctactor {
        public:
            /**
             * Recognises configured objects on the input
             *
             * @param image the bytes that make up the input image
             * @returns the area of the image that has been recognised
             */
            std::vector<cv::Mat> extract(const cv::Mat &&image);
        };

    }
}
