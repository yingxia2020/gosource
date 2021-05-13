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


#include "benchmarks.h"
#include "file_wrapper.h"
#include "options.h"
#include <iostream>

using namespace std;

int main(int argc, char* argv[])
{
    auto options = options::parse(argc, argv);

    //int ret = system("rm -f /tmp/output_deflated_*");
    //ret     = system("rm -f /tmp/output_inflated_*");

    benchmarks benchmarks;

    // adding the benchmark for each files and libary/level combination
    for (const auto& path : options.files)
    {
        auto compression   = benchmark_info::Method::Compression;
        auto decompression = benchmark_info::Method::Decompression;
        auto isal_static   = benchmark_info::Library::ISAL_STATIC;
        auto isal_semidyn  = benchmark_info::Library::ISAL_SEMIDYN;
        auto zlib          = benchmark_info::Library::ZLIB;

        //benchmarks.add_benchmark({compression, isal_static, 0, path});
        benchmarks.add_benchmark({decompression, isal_static, 0, path});

        if (options.isal_semidyn_stateful)
        {
            benchmarks.add_benchmark({compression, isal_semidyn, 0, path});
            benchmarks.add_benchmark({decompression, isal_semidyn, 0, path});
        }
        if (options.isal_semidyn_stateless)
        {
            benchmarks.add_benchmark({compression, isal_semidyn, 1, path});
            benchmarks.add_benchmark({decompression, isal_semidyn, 1, path});
        }

        for (auto level : options.zlib_levels)
        {
            if (level >= 1 && level <= 9)
            {
                benchmarks.add_benchmark({compression, zlib, level, path});
                benchmarks.add_benchmark({decompression, zlib, level, path});
            }
            else
            {
                std::cout << "[Warning] zlib compression level " << level << "will be ignored\n";
            }
        }
    }

    benchmarks.run();
    //int ret1 = system("rm -f /tmp/output_deflated_*");
    //int ret2 = system("rm -f /tmp/output_inflated_*");
    //if (ret1 || ret2)
    //    std::cout << "\n[Warning] Temporary files in /tmp/ have not been deleted\n";
}
