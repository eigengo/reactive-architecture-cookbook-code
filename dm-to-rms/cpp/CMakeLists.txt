cmake_minimum_required(VERSION 3.1.0)
set(CMAKE_CXX_STANDARD 14)

project(cpp)

# Find our dependencies
include(FindProtobuf)
find_package(Protobuf REQUIRED)
find_package(Boost REQUIRED)

set(CMAKE_INCLUDE_CURRENT_DIR TRUE)
set(SRCS demo.cc)
file(GLOB_RECURSE PROTOS ${CMAKE_CURRENT_SOURCE_DIR}/../protocol/*.proto)

#Code Generation
PROTOBUF_GENERATE_CPP(PROTO_SRC PROTO_HEADER ${PROTOS})
INCLUDE_DIRECTORIES(${PROTOBUF_INCLUDE_DIR})

#Executable setup
add_executable(demo ${PROTO_HEADER} ${PROTO_SRC} ${SRCS})
TARGET_LINK_LIBRARIES(demo  ${PROTOBUF_LIBRARY})
