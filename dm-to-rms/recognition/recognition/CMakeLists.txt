cmake_minimum_required(VERSION 3.1.0)
set(CMAKE_CXX_STANDARD 14)

project(recognition)


#Executable setup
add_library(recognition recognition.cc)
