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

#include "benchmark_info.h"
#include <time.h>

std::ostream& operator<<(std::ostream& os, const benchmark_info::Library& lib)
{
    switch (lib)
    {
        case (benchmark_info::Library::ISAL_STATIC): os << std::string("isa-l (static)"); break;
        case (benchmark_info::Library::ISAL_SEMIDYN): os << std::string("isa-l (semi-dyn)"); break;
        case (benchmark_info::Library::ZLIB): os << std::string("zlib"); break;
    }
    return os;
}

std::string benchmark_info::file_name_suffix(benchmark_info info)
{
    srand(time(NULL));
    int randNumb = rand() % 100000 + 1;
    std::string suffix = "_" + std::to_string(std::hash<std::string>()(info.file)+randNumb) + "_";
    switch (info.library)
    {
        case benchmark_info::Library::ISAL_STATIC: suffix += "isa-l-static"; break;
        case benchmark_info::Library::ISAL_SEMIDYN:
            suffix += "isa-l-semidyn_" + std::to_string(info.config);
            break;
        case benchmark_info::Library::ZLIB: suffix += "zlib_" + std::to_string(info.config); break;
    }
    return suffix;
}
