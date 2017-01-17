cmake_minimum_required(VERSION 3.1.0)

project (all)

include(DownloadProject.cmake)
enable_testing()

include_directories(include)
add_subdirectory(recognition)
add_subdirectory(test)

set(CMAKE_CXX_STANDARD 14)

include(FindProtobuf)
find_package(Protobuf REQUIRED)
find_package(Boost REQUIRED)

file(GLOB_RECURSE PROTOS ${CMAKE_CURRENT_SOURCE_DIR}/../protocol/*.proto)

protobuf_generate_cpp(PROTO_SRC PROTO_HEADER ${PROTOS})
set(CMAKE_INCLUDE_CURRENT_DIR TRUE)
include_directories(${PROTOBUF_INCLUDE_DIR})
add_executable(azmq ${PROTO_HEADER} ${PROTO_SRC} azmq.cc)
target_link_libraries(azmq recognition ${PROTOBUF_LIBRARY})
