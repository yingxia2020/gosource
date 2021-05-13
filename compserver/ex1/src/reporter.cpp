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


#include "reporter.h"
#include "file_wrapper.h"
#include "reporter_utils.h"
#include <algorithm>
#include <cmath>
#include <iomanip>
#include <iostream>

using namespace std;

reporter::reporter(int benchmark_count, const vector<benchmark_info>& benchmarks_info)
    : benchmark_count(benchmark_count), benchmark_run(0), benchmarks_info(benchmarks_info)
{
}

// called before the 1st benchmark starts
bool reporter::ReportContext(const benchmark::BenchmarkReporter::Context& context)
{
    cout << "[Info   ] ";
    cout << "Running benchmarks at " << context.mhz_per_cpu << " MHz\n";

    if (context.cpu_scaling_enabled)
        cout << "[Warning] CPU scaling is enabled. Real time measurements will be noisy\n";

    cout << "[Info   ] Benchmarks starting...\n";

    if (!benchmarks_info.empty())
        reporter_utils::display_benchmark_starting(benchmarks_info.front(), 1, benchmark_count);

    return true;
}

// called after each benchmarks are run
// the results are stored and later displayed in the Finalize() method
void reporter::ReportRuns(const vector<benchmark::BenchmarkReporter::Run>& report)
{
    const auto& info = benchmarks_info[benchmark_run];
    benchmark_run++;

    if (report.empty())
        return;

    const auto& run = report.front();

    if (run.error_occurred)
    {
        cout << "[Error  ] ... An error occured:\n";
        cout << "[Error  ] ... \"" << run.error_message << "\"\n";
        cout << "[Error  ] ... This benchmark will be excluded from the final results\n";
    }
    else
    {
        if (info.method == benchmark_info::Method::Compression)
        {
            double ratio = (double)file_wrapper(
                               "/tmp/output_deflated" + benchmark_info::file_name_suffix(info))
                               .size() /
                           file_wrapper(info.file).size();
            ratio = ratio * 100;

            result results = {info,
                              chrono::nanoseconds{(long int)run.GetAdjustedRealTime()},
                              chrono::nanoseconds{(long int)run.GetAdjustedCPUTime()},
                              ratio,
                              chrono::nanoseconds{0},
                              chrono::nanoseconds{0}};

            results_by_file[info.file].push_back(results);
        }
        else
        {
            /*
            results_by_file[info.file].back().decomp_realtime =
                chrono::nanoseconds{(long int)run.GetAdjustedRealTime()};
            results_by_file[info.file].back().decomp_cputime =
                chrono::nanoseconds{(long int)run.GetAdjustedCPUTime()};
                */
            double ratio = (double)file_wrapper(
                               "/tmp/output_inflated" + benchmark_info::file_name_suffix(info))
                               .size() /
                           file_wrapper(info.file).size();
            ratio = ratio * 100;

            result results = {info,
                              chrono::nanoseconds{(long int)run.GetAdjustedRealTime()},
                              chrono::nanoseconds{(long int)run.GetAdjustedCPUTime()},
                              ratio,
                              chrono::nanoseconds{0},
                              chrono::nanoseconds{0}};

            results_by_file[info.file].push_back(results);
        }
    }

    if (benchmark_run < benchmark_count)
        reporter_utils::display_benchmark_starting(
            benchmarks_info[benchmark_run], benchmark_run + 1, benchmark_count);
}

// called after all benchmark run
// display the results for each file, then the average results for each library/level combination
void reporter::Finalize()
{
    // used to store the relative results grouped by library/level combinations
    map<library_and_level, comparative_results> results_by_library;

    // display results for each file
    for (const auto& file_results_pair : results_by_file)
    {
        const auto& file    = file_results_pair.first;
        const auto& results = file_results_pair.second;

        if (results.empty())
            continue;

        cout << "\n" << file << " (";
        reporter_utils::display_size(file_wrapper(file).size());
        cout << ")\n";
        reporter_utils::display_header();

        // the first benchmark for this file is considered the 'base' (it always is isa-l)
        // relative results will be relative to this base
        const auto& base = results.front();

        for (const auto& bm : results)
        {
            // compare the result relative to the base
            comparative_results comparative_result{
                (double)bm.comp_realtime.count() / base.comp_realtime.count(),
                (double)bm.comp_cputime.count() / base.comp_cputime.count(),
                bm.comp_ratio / base.comp_ratio,
                (double)bm.decomp_realtime.count() / base.decomp_realtime.count(),
                (double)bm.decomp_cputime.count() / base.decomp_cputime.count()};

            // accumulate the relative results for this library/level combination
            results_by_library[{bm.info.library, bm.info.config}] += comparative_result;

            reporter_utils::display_file_result_line(bm, comparative_result);
        }
    }

    // compute the average relative results for each libary/level combination
    for (auto& p : results_by_library)
    {
        p.second /= p.second.aggregate_count;
    }

    // display the average results for each library/level combination
    cout << "\nAverage results:\n";
    reporter_utils::display_header();
    for (const auto& bm : results_by_library)
    {
        reporter_utils::display_average_result_line(bm.first.library, bm.first.level, bm.second);
    }
}
