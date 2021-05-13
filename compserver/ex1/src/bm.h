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


#ifndef EXS_BM_H
#define EXS_BM_H

#include "benchmark_info.h"
#include "file_wrapper.h"
#include <benchmark/benchmark.h>
#include <chrono>
#include <memory>

// the base class for each library to be benchmark
//
// the run() method will be called by google benchmark, and will run the overrided iteration()
// method enough times to get a good time measurement
//
// the iteration() method performs a compression once, and return how long it took
//
// the has_level() method must be overrided and set to return 'true' if the library accepts a
// compression level
//
// implementations of this abstract class are in the bm_* files
class bm
{
  public:
    using raw_duration = std::chrono::steady_clock::duration;

    virtual ~bm() = default;

    virtual void
    run_deflate(benchmark::State& state, std::string in_file_name, std::string out_file_name) final;

    virtual void
    run_inflate(benchmark::State& state, std::string in_file_name, std::string out_file_name) final;

    virtual std::string version() = 0;

  protected:
    virtual raw_duration
    iter_deflate(file_wrapper* in_file, file_wrapper* out_file, int config) = 0;
    virtual raw_duration iter_inflate(file_wrapper* in_file, file_wrapper* out_file) = 0;
    virtual bool has_config() const;
};

#endif
