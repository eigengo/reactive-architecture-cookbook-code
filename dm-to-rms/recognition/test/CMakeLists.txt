cmake_minimum_required(VERSION 3.1.0)
set(CMAKE_CXX_STANDARD 14)

project(main-test)

include_directories("${gtest_SOURCE_DIR}/include" "${gmock_SOURCE_DIR}/include")

find_package(OpenCV 3 REQUIRED)
include_directories(${OpenCV_INCLUDE_DIRS})
download_project(PROJ       rapidcheck
        GIT_REPOSITORY      https://github.com/emil-e/rapidcheck
        GIT_TAG             master
        ${UPDATE_DISCONNECTED_IF_AVAILABLE}
        )

add_subdirectory(${rapidcheck_SOURCE_DIR} ${rapidcheck_BINARY_DIR})
include_directories(${rapidcheck_SOURCE_DIR}/include ${rapidcheck_SOURCE_DIR}/extras/gtest/include)
include_directories(${PROTOBUF_INCLUDE_DIR})
include_directories(${CMAKE_CURRENT_BINARY_DIR}/..)

message(STATUS "************ ${PROTO_HEADER} ${PROTO_SRC}")

set(SRCS main-test.cc)

add_executable(main-test ${SRCS} ${PROTO_HEADER} ${PROTO_SRC})
target_link_libraries(main-test rapidcheck gtest gmock_main ${PROTOBUF_LIBRARY} recognition)
add_test(all main-test)
