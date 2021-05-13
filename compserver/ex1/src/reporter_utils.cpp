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


#include "reporter_utils.h"
#include "file_wrapper.h"
#include <cmath>
#include <iomanip>
#include <iostream>
#include <vector>

using namespace std;

void reporter_utils::display_size(long size)
{
    if (size < 1000)
    {
        cout << size << " B";
        return;
    }

    vector<string> prefixes = {"k", "M", "G"};

    int    exp  = (int)(log(size) / log(1000));
    string unit = prefixes[exp - 1] + "iB";

    cout << fixed << setprecision(1) << size / pow(1000, exp) << " " << unit;
}

void reporter_utils::display_benchmark_starting(
    const benchmark_info& info, int bm_run, int bm_count)
{
    cout << "[Info   ] ";
    cout << "(" << setw((int)log10(bm_count) + 1) << right << bm_run << "/" << bm_count << ") ";

    if (info.method == benchmark_info::Method::Compression)
    {
        cout << "Compressing " << info.file << " (";
        display_size(file_wrapper(info.file).size());
        cout << ") with " << info.library;
    }
    else
    {
        cout << "Decompressing ";
    }

    cout << "...\n";
}

void reporter_utils::display_header()
{
    // clang-format off
    cout << "---------------------------------------------------------------------------------------------------------------------------------------------\n";
    cout << " Library         | Config    | Compression                                                      | Decompression                               \n";
    cout << "                 |           | Ratio              | Real Time            | CPU Time             | Real Time            | CPU Time             \n";
    cout << "-----------------|-----------|--------------------|----------------------|----------------------|----------------------|----------------------\n";
    // clang-format on
}

std::string duration_to_string(double dur)
{
    if (dur > 100000)
        return std::to_string((int64_t)(dur / 1000)) + " us";
    else
        return std::to_string((int64_t)dur) + " ns";
}

void reporter_utils::display_file_result_line(
    const result& result, const comparative_results& comparative_result)
{
    // library name
    cout << setw(16) << left << result.info.library << " | ";

    // library config
    if (result.info.library == benchmark_info::Library::ISAL_STATIC)
        cout << setw(9) << right << "-"
             << " | ";
    else if (result.info.library == benchmark_info::Library::ISAL_SEMIDYN)
        cout << setw(9) << right << (result.info.config ? "stateless" : "stateful") << " | ";
    else
        cout << setw(9) << right << ("level " + std::to_string(result.info.config)) << " | ";

    // compression ratio
    cout << setw(7) << right << fixed << setprecision(2) << result.comp_ratio << " % (x" << setw(5)
         << fixed << setprecision(2) << right << comparative_result.comp_ratio << ") | ";

    // compression real-time
    cout << setw(10) << right << duration_to_string(result.comp_realtime.count()) << " (x"
         << setw(6) << fixed << setprecision(2) << right << comparative_result.comp_realtime
         << ") | ";

    // compression cpu-time
    cout << setw(10) << right << duration_to_string(result.comp_cputime.count()) << " (x" << setw(6)
         << fixed << setprecision(2) << right << comparative_result.comp_cputime << ") | ";

    // decompression real-time
    cout << setw(10) << right << duration_to_string(result.decomp_realtime.count()) << " (x"
         << setw(6) << fixed << setprecision(2) << right << comparative_result.decomp_realtime
         << ") | ";

    // decompression cpu-time
    cout << setw(10) << right << duration_to_string(result.decomp_cputime.count()) << " (x"
         << setw(6) << fixed << setprecision(2) << right << comparative_result.decomp_cputime
         << ")";

    cout << "\n";
}

void reporter_utils::display_average_result_line(
    benchmark_info::Library library, int config, const comparative_results& result)
{
    // library name
    cout << setw(16) << left << library << " | ";

    // library config
    if (library == benchmark_info::Library::ISAL_STATIC)
        cout << setw(9) << right << "-"
             << " | ";
    else if (library == benchmark_info::Library::ISAL_SEMIDYN)
        cout << setw(9) << right << (config ? "stateless" : "stateful") << " | ";
    else
        cout << setw(9) << right << ("level " + std::to_string(config)) << " | ";

    cout << "x" << setw(17) << fixed << setprecision(2) << right << fixed << setprecision(2)
         << result.comp_ratio << " | ";
    cout << "x" << setw(19) << fixed << setprecision(2) << right << result.comp_realtime << " | ";
    cout << "x" << setw(19) << fixed << setprecision(2) << right << result.comp_cputime << " | ";
    cout << "x" << setw(19) << fixed << setprecision(2) << right << result.decomp_realtime << " | ";
    cout << "x" << setw(19) << fixed << setprecision(2) << right << result.decomp_cputime;
    cout << "\n";
}
