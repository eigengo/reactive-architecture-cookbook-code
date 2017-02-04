cmake_minimum_required(VERSION 3.1.0)
set(CMAKE_CXX_STANDARD 14)

project (all)

list(APPEND CMAKE_MODULE_PATH ${CMAKE_CURRENT_SOURCE_DIR}/modules)

include(DownloadProject.cmake)
include(FindProtobuf)
find_package(Protobuf REQUIRED)
find_package(Boost COMPONENTS system regex program_options REQUIRED)
find_package(LibRDKafka REQUIRED)
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
add_subdirectory(faceextract)
add_subdirectory(faceextract-test)
add_subdirectory(test)

add_executable(main ${PROTO_HEADER} ${PROTO_SRC} main.cc)
target_link_libraries(main faceextract pthread ${PROTOBUF_LIBRARY} ${Boost_LIBRARIES} ${LibRDKafka_LIBRARIES})
