cmake_minimum_required(VERSION 3.1.0)
set(CMAKE_CXX_STANDARD 14)

project (all)

include(DownloadProject.cmake)
include(FindProtobuf)
find_package(Protobuf REQUIRED)
find_package(Boost COMPONENTS system regex program_options REQUIRED)
enable_testing()

download_project(PROJ       googletest
        GIT_REPOSITORY      https://github.com/google/googletest.git
        GIT_TAG             master
        ${UPDATE_DISCONNECTED_IF_AVAILABLE}
        )

set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
add_subdirectory(${googletest_SOURCE_DIR} ${googletest_BINARY_DIR})
include_directories("${gtest_SOURCE_DIR}/include" "${gmock_SOURCE_DIR}/include")

file(GLOB_RECURSE PROTOS ${CMAKE_CURRENT_SOURCE_DIR}/../protocol/*.proto)
protobuf_generate_cpp(PROTO_SRC PROTO_HEADER ${PROTOS})
set(CMAKE_INCLUDE_CURRENT_DIR TRUE)
include_directories(${PROTOBUF_INCLUDE_DIR})

include_directories(include)
add_subdirectory(recognition)
add_subdirectory(recognition-test)
add_subdirectory(test)

download_project(PROJ       azmq
        GIT_REPOSITORY      https://github.com/zeromq/azmq.git
        GIT_TAG             master
        ${UPDATE_DISCONNECTED_IF_AVAILABLE}
        )

list(APPEND CMAKE_MODULE_PATH "${azmq_SOURCE_DIR}/config/")
find_package(ZeroMQ 4.0 REQUIRED)
include_directories("${azmq_SOURCE_DIR}")

add_executable(main ${PROTO_HEADER} ${PROTO_SRC} main.cc)
target_link_libraries(main recognition ${PROTOBUF_LIBRARY} ${Boost_LIBRARIES} ${ZeroMQ_LIBRARIES})
