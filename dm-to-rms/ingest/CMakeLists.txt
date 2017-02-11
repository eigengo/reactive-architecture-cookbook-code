cmake_minimum_required(VERSION 3.1.0)
set(CMAKE_CXX_STANDARD 14)

project (all)

list(APPEND CMAKE_MODULE_PATH ${CMAKE_CURRENT_SOURCE_DIR}/modules)
set(OPENSSL_ROOT_DIR /usr/local/opt/openssl)

include(DownloadProject.cmake)
include(FindProtobuf)
find_package(LibRDKafka REQUIRED)
find_package(Protobuf REQUIRED)
find_package(ngHTTP2 REQUIRED)
find_package(OpenSSL REQUIRED)
find_package(Boost COMPONENTS system regex program_options REQUIRED)
enable_testing()

# This is OSX & Homebrew hack
if (APPLE)
    set(OPENSSL_LIBRARIES ${OPENSSL_ROOT_DIR}/lib/libssl.dylib ${OPENSSL_ROOT_DIR}/lib/libcrypto.dylib )
endif()

download_project(PROJ       googletest
        GIT_REPOSITORY      https://github.com/google/googletest.git
        GIT_TAG             master
        ${UPDATE_DISCONNECTED_IF_AVAILABLE})

download_project(PROJ       easyloggingpp
        GIT_REPOSITORY      https://github.com/muflihun/easyloggingpp
        GIT_TAG             master
        ${UPDATE_DISCONNECTED_IF_AVAILABLE}
        )

set(easyloggingpp_SRC ${easyloggingpp_SOURCE_DIR}/src/easylogging++.cc)
include_directories(${easyloggingpp_SOURCE_DIR}/src)

set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
add_subdirectory(${googletest_SOURCE_DIR} ${googletest_BINARY_DIR})
include_directories("${gtest_SOURCE_DIR}/include" "${gmock_SOURCE_DIR}/include" "${libkafka-asio_SOURCE_DIR}/lib")

file(GLOB_RECURSE PROTOS ${CMAKE_CURRENT_SOURCE_DIR}/../protocol/*.proto)
protobuf_generate_cpp(PROTO_SRC PROTO_HEADER ${PROTOS})
set(CMAKE_INCLUDE_CURRENT_DIR TRUE)
include_directories(${PROTOBUF_INCLUDE_DIR})
include_directories(${NGHTTP2_INCLUDE_DIR} ${OPENSSL_INCLUDE_DIR} ${SNAPPY_INCLUDE_DIR})

include_directories(include)
add_subdirectory(ingest)
add_subdirectory(ingest-test)
add_subdirectory(test)

message(STATUS "*********** ${LibRDKafka_LIBRARIES} : ${SNAPPY_INCLUDE_DIR}")

add_executable(main ${PROTO_HEADER} ${PROTO_SRC} main.cc)
target_link_libraries(main ingest pthread ${PROTOBUF_LIBRARY} ${LibRDKafka_LIBRARIES} ${Boost_LIBRARIES} ${NGHTTP2_LIBRARY} ${NGHTTP2_ASIO_LIBRARY} ${OPENSSL_LIBRARIES})
