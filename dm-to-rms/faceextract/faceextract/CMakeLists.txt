cmake_minimum_required(VERSION 3.1.0)
set(CMAKE_CXX_STANDARD 14)

project(faceextractor)

find_package(OpenCV 3 REQUIRED)
include_directories(${OpenCV_INCLUDE_DIRS})

#Executable setup
add_library(faceextractor faceextractor.cc)
target_link_libraries(faceextractor ${OpenCV_LIBS})
