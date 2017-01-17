cmake_minimum_required(VERSION 3.1.0)
set(CMAKE_CXX_STANDARD 14)

project (all)

include(DownloadProject.cmake)
include(FindProtobuf)
find_package(Protobuf REQUIRED)
find_package(Boost REQUIRED)
enable_testing()

include_directories(include)
add_subdirectory(recognition)
add_subdirectory(test)

file(GLOB_RECURSE PROTOS ${CMAKE_CURRENT_SOURCE_DIR}/../protocol/*.proto)
protobuf_generate_cpp(PROTO_SRC PROTO_HEADER ${PROTOS})
set(CMAKE_INCLUDE_CURRENT_DIR TRUE)
include_directories(${PROTOBUF_INCLUDE_DIR})

download_project(PROJ       azmq
        GIT_REPOSITORY      https://github.com/zeromq/azmq.git
        GIT_TAG             master
        ${UPDATE_DISCONNECTED_IF_AVAILABLE}
        )

#
# set ${CMAKE_BUILD_DIR}/azmq-src/CMakeList.txt:36 to
# list(APPEND CMAKE_MODULE_PATH "${CMAKE_SOURCE_DIR}/config")
#

list(APPEND CMAKE_MODULE_PATH "${azmq_SOURCE_DIR}/config/")
include(${azmq_SOURCE_DIR}/config/FindZeroMQ.cmake)
add_subdirectory(${azmq_SOURCE_DIR} ${azmq_BINARY_DIR})
include_directories("${azmq_SOURCE_DIR}/include")

add_executable(main ${PROTO_HEADER} ${PROTO_SRC} main.cc)
target_link_libraries(main recognition azmq ${PROTOBUF_LIBRARY})
