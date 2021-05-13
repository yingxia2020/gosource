/*
Copyright (c) <2016>, Intel Corporation

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice,
this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.
* Neither the name of Intel Corporation nor the names of its contributors
may be used to endorse or promote products derived from this software
without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/


#include "options.h"

#include <boost/algorithm/string.hpp>
#include <boost/filesystem.hpp>
#include <boost/program_options.hpp>
#include <iostream>

std::vector<int> parse_compression_levels(std::string str)
{
    std::vector<std::string> levels_str;
    boost::split(levels_str, str, boost::is_any_of(", "));

    std::vector<int> levels;
    for (const auto& level_str : levels_str)
        levels.push_back(std::stoi(level_str));

    return levels;
}

void list_directory(std::string path, std::vector<std::string>& files)
{
    using namespace boost::filesystem;

    if (!exists(path) || !is_directory(path))
        return;

    for (auto& entry : recursive_directory_iterator(path))
    {
        if (entry.status().type() == file_type::regular_file)
            files.push_back(entry.path().string());
    }
}

options options::parse(int argc, char* argv[])
{
    using namespace boost::program_options;
    using std::vector;
    using std::string;
    using std::cout;

    options_description desc("Usage: ./ex1 [--help] [--folder <path>]... [--file <path>]... ");

    desc.add_options()("help", "display this message");
    desc.add_options()(
        "file", value<vector<string>>()->value_name("path"), "use the file at 'path'");
    desc.add_options()(
        "folder", value<vector<string>>()->value_name("path"), "use all the files in 'path'");
    desc.add_options()(
        "zlib-levels",
        value<string>()->value_name("n,..."),
        "coma-separated list of compression level [1-9]");
    desc.add_options()(
        "semidyn-config",
        value<string>()->value_name("flag,..."),
        "coma-separated list of flags "
        "for the semi-dynamic "
        "compression ('stateful', "
        "'stateless') [stateful]");

    variables_map vm;

    try
    {
        store(parse_command_line(argc, argv, desc), vm);
    }
    catch (unknown_option&)
    {
    }

    notify(vm);

    if (vm.count("help") || (!vm.count("folder") && !vm.count("file")))
    {
        cout << desc << "\n";
        exit(0);
    }

    vector<string> files = vm.count("file") ? vm["file"].as<vector<string>>() : vector<string>{};
    vector<string> folders =
        vm.count("folder") ? vm["folder"].as<vector<string>>() : vector<string>{};
    for (const auto& folder : folders)
    {
        list_directory(folder, files);
    }

    bool stateful = vm.count("semidyn-config") == 0 ||
                    vm["semidyn-config"].as<std::string>().find("ful") != std::string::npos;
    bool stateless = vm.count("semidyn-config") != 0 &&
                     vm["semidyn-config"].as<std::string>().find("less") != std::string::npos;

    return options{
        files,
        parse_compression_levels(vm.count("zlib-levels") ? vm["zlib-levels"].as<string>() : "6"),
        stateful,
        stateless};
}
