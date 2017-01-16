#include <string>
#include <memory>
#include <vector>

namespace com {
    namespace reactivearchitecturecookbook {

        class Recogniser {
        public:
            /**
             * Recognises configured objects on the input, which should be an image
             * @param image the bytes that make up the input image
             * @returns the area of the image that has been recognised
             */
            std::vector<uint8_t> recognise(std::vector<uint8_t> &image);
        };

    }
}
