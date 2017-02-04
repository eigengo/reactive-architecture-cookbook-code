cmake_minimum_required(VERSION 3.1.0)
set(CMAKE_CXX_STANDARD 14)

project(recognition-test)

include_directories("${gtest_SOURCE_DIR}/include" "${gmock_SOURCE_DIR}/include")

find_package(OpenCV 3 REQUIRED)
include_directories(${OpenCV_INCLUDE_DIRS})

set(SRCS faceextractor-test.cc)

add_executable(recognition-test ${SRCS})
target_link_libraries(recognition-test gtest gmock_main ${OpenCV_LIBS} recognition)
add_test(all recognition-test)
