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


#include "bm.h"

void bm::run_deflate(benchmark::State& state, std::string in_file_name, std::string out_file_name)
{
    file_wrapper out_file(out_file_name);
    if (!out_file.open(file_wrapper::open_mode::W))
    {
        state.SkipWithError("Cannot open output file for writing.");
    }

    file_wrapper in_file(in_file_name);
    if (!in_file.open(file_wrapper::open_mode::R))
    {
        state.SkipWithError("Cannot open input file for writing.");
    }

    int config = 0;
    if (has_config())
        config = state.range(0);

    while (state.KeepRunning())
    {
        in_file.seek_to_start();
        out_file.truncate();

        auto duration = iter_deflate(&in_file, &out_file, config);
        if (duration == raw_duration{0})
        {
            state.SkipWithError("Error running deflate function");
            break;
        }
        state.SetIterationTime(
            std::chrono::duration_cast<std::chrono::duration<double>>(duration).count());
    }
}

void bm::run_inflate(benchmark::State& state, std::string in_file_name, std::string out_file_name)
{
    file_wrapper out_file(out_file_name);
    if (!out_file.open(file_wrapper::open_mode::W))
    {
        state.SkipWithError("Cannot open output file for writing.");
    }

    file_wrapper in_file(in_file_name);
    if (!in_file.open(file_wrapper::open_mode::R))
    {
        state.SkipWithError("Cannot open input file for writing.");
    }

    while (state.KeepRunning())
    {
        in_file.seek_to_start();
        out_file.truncate();

        auto duration = iter_inflate(&in_file, &out_file);
        if (duration == raw_duration{0})
        {
            state.SkipWithError("Error running inflate function");
            break;
        }
        state.SetIterationTime(
            std::chrono::duration_cast<std::chrono::duration<double>>(duration).count());
    }
}

bool bm::has_config() const
{
    return true;
}
