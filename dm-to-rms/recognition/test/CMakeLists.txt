cmake_minimum_required(VERSION 3.1.0)
set(CMAKE_CXX_STANDARD 14)

project(test)

download_project(PROJ                googletest
                 GIT_REPOSITORY      https://github.com/google/googletest.git
                 GIT_TAG             master
                 ${UPDATE_DISCONNECTED_IF_AVAILABLE}
)

set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
add_subdirectory(${googletest_SOURCE_DIR} ${googletest_BINARY_DIR})
include_directories("${gtest_SOURCE_DIR}/include" "${gmock_SOURCE_DIR}/include")

find_package(OpenCV 3 REQUIRED)
include_directories(${OpenCV_INCLUDE_DIRS})

set(SRCS recogniser-test.cc)

add_executable(recognition-test ${SRCS})
target_link_libraries(recognition-test gtest gmock_main ${OpenCV_LIBS} recognition)
add_test(all recognition-test)
