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


#include "bm_isal.h"
#include <isa-l.h>

#define BUF_SIZE 262144

std::string bm_isal::version()
{
    return std::to_string(ISAL_MAJOR_VERSION) + "." + std::to_string(ISAL_MINOR_VERSION) + "." +
           std::to_string(ISAL_PATCH_VERSION);
}

bm::raw_duration
bm_isal::iter_deflate(file_wrapper* in_file, file_wrapper* out_file, int /*config*/)
{
    raw_duration duration{};

    struct isal_zstream stream;

    uint8_t input_buffer[BUF_SIZE];
    uint8_t output_buffer[BUF_SIZE];

    isal_deflate_init(&stream);
    stream.end_of_stream = 0;
    stream.flush         = NO_FLUSH;

    do
    {
        stream.avail_in      = static_cast<uint32_t>(in_file->read(input_buffer, BUF_SIZE));
        stream.end_of_stream = static_cast<uint32_t>(in_file->eof());
        stream.next_in       = input_buffer;
        do
        {
            stream.avail_out = BUF_SIZE;
            stream.next_out  = output_buffer;

            auto begin = std::chrono::steady_clock::now();
            isal_deflate(&stream);
            auto end = std::chrono::steady_clock::now();
            duration += (end - begin);

            out_file->write(output_buffer, BUF_SIZE - stream.avail_out);
        } while (stream.avail_out == 0);
    } while (stream.internal_state.state != ZSTATE_END);

    return duration;
}

bm::raw_duration bm_isal::iter_inflate(file_wrapper* in_file, file_wrapper* out_file)
{
    raw_duration duration{};

    int                  ret;
    int                  eof;
    struct inflate_state stream;

    uint8_t input_buffer[BUF_SIZE];
    uint8_t output_buffer[BUF_SIZE];

    isal_inflate_init(&stream);

    stream.avail_in = 0;
    stream.next_in  = nullptr;

    do
    {
        stream.avail_in = static_cast<uint32_t>(in_file->read(input_buffer, BUF_SIZE));
        eof             = in_file->eof();
        stream.next_in  = input_buffer;
        do
        {
            stream.avail_out = BUF_SIZE;
            stream.next_out  = output_buffer;

            auto begin = std::chrono::steady_clock::now();
            ret        = isal_inflate(&stream);
            auto end   = std::chrono::steady_clock::now();
            duration += (end - begin);

            out_file->write(output_buffer, BUF_SIZE - stream.avail_out);
        } while (stream.avail_out == 0);
    } while (ret != ISAL_END_INPUT && eof == 0);

    return duration;
}

bool bm_isal::has_config() const
{
    return false;
}
