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


#ifndef REPORTER_H
#define REPORTER_H

#include "benchmark_info.h"
#include <benchmark/reporter.h>
#include <chrono>
#include <map>
#include <unordered_map>

// store the results of a benchmark
struct result
{
    benchmark_info           info;
    std::chrono::nanoseconds comp_realtime;
    std::chrono::nanoseconds comp_cputime;
    double                   comp_ratio;
    std::chrono::nanoseconds decomp_realtime;
    std::chrono::nanoseconds decomp_cputime;
};

// a convenience struct combining the library and the compression level used by a benchmark
struct library_and_level
{
    benchmark_info::Library library;
    int                     level;

    bool operator<(const library_and_level& o) const
    {
        if (library == o.library)
            return level < o.level;
        return library < o.library;
    }
};

// store the relative results of a benchmark compared to a base results
struct comparative_results
{
    double comp_realtime   = 0;
    double comp_cputime    = 0;
    double comp_ratio      = 0;
    double decomp_realtime = 0;
    double decomp_cputime  = 0;
    int    aggregate_count = 0;

    comparative_results operator+(const comparative_results& o)
    {
        return {comp_realtime + o.comp_realtime,
                comp_cputime + o.comp_cputime,
                comp_ratio + o.comp_ratio,
                decomp_realtime + o.decomp_realtime,
                decomp_cputime + o.decomp_cputime};
    }

    comparative_results& operator+=(const comparative_results& o)
    {
        comp_realtime += o.comp_realtime;
        comp_cputime += o.comp_cputime;
        comp_ratio += o.comp_ratio;
        decomp_realtime += o.decomp_realtime;
        decomp_cputime += o.decomp_cputime;
        aggregate_count++;
        return *this;
    }

    comparative_results operator/(double d)
    {
        return {comp_realtime / d,
                comp_cputime / d,
                comp_ratio / d,
                decomp_realtime / d,
                decomp_cputime / d};
    }

    comparative_results& operator/=(double d)
    {
        comp_realtime /= d;
        comp_cputime /= d;
        comp_ratio /= d;
        decomp_realtime /= d;
        decomp_cputime /= d;
        return *this;
    }
};

// display the benchmarks info and results
class reporter : public benchmark::BenchmarkReporter
{
  public:
    reporter(int benchmark_count, const std::vector<benchmark_info>& benchmarks_info);

    // called before the 1st benchmark starts
    virtual bool ReportContext(const Context& context) override;
    // called after each benchmark
    virtual void ReportRuns(const std::vector<Run>& report) override;
    // called after all benchmark run
    virtual void Finalize() override;

  private:
    int                                benchmark_count; // the total number of benchmark to run
    int                                benchmark_run;   // the number of benchmark run so far
    const std::vector<benchmark_info>& benchmarks_info; // the benchmarks to run
    std::unordered_map<std::string, std::vector<result>> results_by_file;
};

#endif
