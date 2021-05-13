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
#include "reporter.h"
#include <functional>
#include <iostream>

benchmarks::benchmarks() : benchmark_count(0)

{
    isal_static  = std::make_unique<bm_isal>();
    isal_semidyn = std::make_unique<bm_isal_semidyn>();
    zlib         = std::make_unique<bm_zlib>();
}

void benchmarks::add_benchmark(benchmark_info info)
{
    std::string file_name_suffix = benchmark_info::file_name_suffix(info);

    auto function_compress =
        [&, info](benchmark::State& state, std::string in_file_name, std::string file_name_suffix) {
            std::string out_file_name = "/tmp/output_deflated" + file_name_suffix;
            switch (info.library)
            {
                case benchmark_info::Library::ISAL_STATIC:
                    isal_static->run_deflate(state, in_file_name, out_file_name);
                    break;
                case benchmark_info::Library::ISAL_SEMIDYN:
                    isal_semidyn->run_deflate(state, in_file_name, out_file_name);
                    break;
                case benchmark_info::Library::ZLIB:
                    zlib->run_deflate(state, in_file_name, out_file_name);
                    break;
            }
        };

    auto function_decompress = 
	[&, info](benchmark::State& state,  std::string in_file_name, std::string file_name_suffix) {
        std::string out_file_name = "/tmp/output_inflated" + file_name_suffix;
        // std::string in_file_name  = "/tmp/output_deflated" + file_name_suffix;

        switch (info.library)
        {
            case benchmark_info::Library::ISAL_STATIC:
                isal_static->run_inflate(state, in_file_name, out_file_name);
                break;
            case benchmark_info::Library::ISAL_SEMIDYN:
                isal_semidyn->run_inflate(state, in_file_name, out_file_name);
                break;
            case benchmark_info::Library::ZLIB:
                zlib->run_inflate(state, in_file_name, out_file_name);
                break;
        }
    };

    auto registration =
        info.method == benchmark_info::Method::Compression
            ? benchmark::RegisterBenchmark("", function_compress, info.file, file_name_suffix)
            : benchmark::RegisterBenchmark("", function_decompress, info.file, file_name_suffix);
    registration->UseManualTime();
    registration->Unit(benchmark::kNanosecond);
    registration->Arg(info.config);

    benchmark_count++;
    benchmarks_info.push_back(info);
}

void benchmarks::run()
{
    int a = 0;
    benchmark::Initialize(&a, nullptr);

    std::cout << "[Info   ] Using isa-l      " << isal_static->version() << "\n";
    std::cout << "[Info   ] Using zlib       " << zlib->version() << "\n";

    auto console_reporter = std::make_unique<reporter>(benchmark_count, benchmarks_info);
    benchmark::RunSpecifiedBenchmarks(console_reporter.get());
}
