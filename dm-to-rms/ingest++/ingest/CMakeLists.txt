cmake_minimum_required(VERSION 3.1.0)
set(CMAKE_CXX_STANDARD 14)

project(ingest)

#Executable setup
add_library(ingest ingest.cc)
target_link_libraries(ingest)
