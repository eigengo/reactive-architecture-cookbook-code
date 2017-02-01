cmake_minimum_required(VERSION 3.1.0)
set(CMAKE_CXX_STANDARD 14)

project(recognition-test)

include_directories("${gtest_SOURCE_DIR}/include" "${gmock_SOURCE_DIR}/include")

set(SRCS ingest-test.cc)

add_executable(ingest-test ${SRCS})
target_link_libraries(ingest-test gtest gmock_main ingest)
add_test(all ingest-test)
