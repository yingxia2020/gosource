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


#include "file_wrapper.h"
#include <sys/stat.h>
#include <unistd.h>

using std::string;

file_wrapper::file_wrapper(string path) : m_file(nullptr), m_path(std::move(path))
{
}

file_wrapper::~file_wrapper()
{
    if (m_file != nullptr)
        fclose(m_file);
}

bool file_wrapper::open(open_mode mode)
{
    m_file = fopen(m_path.c_str(), mode == open_mode::R ? "r" : "w");
    return m_file != nullptr;
}

bool file_wrapper::is_open() const
{
    return m_file != nullptr;
}

bool file_wrapper::seek_to_start()
{
    return fseek(m_file, 0, SEEK_SET) == 0;
}

bool file_wrapper::truncate()
{
    return seek_to_start() && ftruncate(fileno(m_file), 0) == 0;
}

const string& file_wrapper::path() const
{
    return m_path;
}

long file_wrapper::size() const
{
    struct stat st;
    stat(m_path.c_str(), &st);
    return st.st_size;
}

size_t file_wrapper::read(uint8_t* dst_buffer, size_t size)
{
    if (m_file == nullptr)
        return 0;
    return fread(dst_buffer, 1, size, m_file);
}

size_t file_wrapper::write(uint8_t* src_buffer, size_t size)
{
    if (m_file == nullptr)
        return 0;
    return fwrite(src_buffer, 1, size, m_file);
}

int file_wrapper::eof() const
{
    return feof(m_file);
}

int file_wrapper::error() const
{
    return ferror(m_file);
}
