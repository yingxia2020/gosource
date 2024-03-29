project(ex1)
cmake_minimum_required(VERSION 3.1)

SET(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++14")

include_directories(SYSTEM ./libs/usr/include)
link_directories(./libs/usr/lib)

link_libraries(isal)
link_libraries(benchmark)
link_libraries(boost_program_options)
link_libraries(boost_system)
link_libraries(boost_filesystem)
link_libraries(z)

# benchmark library requires pthread
set(THREADS_PREFER_PTHREAD_FLAG ON)
find_package(Threads REQUIRED)
link_libraries(Threads::Threads)

# export compile commands file to build directory
SET(CMAKE_EXPORT_COMPILE_COMMANDS ON)

aux_source_directory(src SRC_LIST)
add_executable(${PROJECT_NAME} ${SRC_LIST})
