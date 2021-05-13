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


#include "bm_isal_semidyn.h"
#include <cstring>
#include <isa-l.h>

#define SEGMENT_SIZE (512 * 1024)
#define SAMPLE_SIZE (32 * 1024)
#define INFLATE_BUF_SIZE 262144

std::string bm_isal_semidyn::version()
{
    return std::to_string(ISAL_MAJOR_VERSION) + "." + std::to_string(ISAL_MINOR_VERSION) + "." +
           std::to_string(ISAL_PATCH_VERSION);
}

bm::raw_duration
bm_isal_semidyn::iter_deflate(file_wrapper* in_file, file_wrapper* out_file, int config)
{
    raw_duration duration{};

    bool stateful = (config == 0);

    struct isal_zstream        stream;
    struct isal_huff_histogram histogram;
    struct isal_hufftables     hufftable;

    long     in_file_size = in_file->size();
    uint8_t* input_buffer = new (std::nothrow) uint8_t[in_file_size];

    if (input_buffer == nullptr)
        return raw_duration{0};

    long     out_buffer_size = std::max((int)(in_file_size * 1.30), 4 * 1024);
    uint8_t* output_buffer   = new (std::nothrow) uint8_t[out_buffer_size];

    if (output_buffer == nullptr)
        return raw_duration{0};

    stream.avail_in = static_cast<uint32_t>(in_file->read(input_buffer, in_file_size));
    if (stream.avail_in != in_file_size)
        return raw_duration{0};

    int segment_size = SEGMENT_SIZE;
    int sample_size  = SAMPLE_SIZE;
    int hist_size    = sample_size > segment_size ? segment_size : sample_size;

    if (stateful)
        isal_deflate_init(&stream);
    else
        isal_deflate_stateless_init(&stream);

    stream.end_of_stream = 0;
    stream.flush         = stateful ? SYNC_FLUSH : FULL_FLUSH;
    stream.next_in       = input_buffer;
    stream.next_out      = output_buffer;
    if (stateful)
        stream.avail_out = out_buffer_size;
    int remaining        = in_file_size;
    int chunk_size       = segment_size;

    while (remaining > 0)
    {
        auto step = std::chrono::steady_clock::now();
        memset(&histogram, 0, sizeof(struct isal_huff_histogram));
        duration += std::chrono::steady_clock::now() - step;

        if (remaining < segment_size * 2)
        {
            chunk_size           = remaining;
            stream.end_of_stream = 1;
        }

        step         = std::chrono::steady_clock::now();
        int hist_rem = (hist_size > chunk_size) ? chunk_size : hist_size;
        isal_update_histogram(stream.next_in, hist_rem, &histogram);
        if (hist_rem == chunk_size)
            isal_create_hufftables_subset(&hufftable, &histogram);
        else
            isal_create_hufftables(&hufftable, &histogram);
        duration += std::chrono::steady_clock::now() - step;

        stream.avail_in = chunk_size;
        if (!stateful)
            stream.avail_out = chunk_size + 8 * (1 + (chunk_size >> 16));

        stream.hufftables = &hufftable;
        remaining -= chunk_size;
        step = std::chrono::steady_clock::now();
        if (stateful)
            isal_deflate(&stream);
        else
            isal_deflate_stateless(&stream);
        duration += std::chrono::steady_clock::now() - step;

        if (stateful)
        {
            if (stream.internal_state.state != ZSTATE_NEW_HDR)
                break;
        }
        else
        {
            if (stream.avail_in != 0)
                break;
        }
    }

    if (stream.avail_in != 0)
        return raw_duration{0};

    out_file->write(output_buffer, stream.total_out);

    delete[] input_buffer;
    delete[] output_buffer;

    return duration;
}

bm::raw_duration bm_isal_semidyn::iter_inflate(file_wrapper* in_file, file_wrapper* out_file)
{
    raw_duration duration{};

    int                  ret;
    int                  eof;
    struct inflate_state stream;

    uint8_t input_buffer[INFLATE_BUF_SIZE];
    uint8_t output_buffer[INFLATE_BUF_SIZE];

    isal_inflate_init(&stream);

    stream.avail_in = 0;
    stream.next_in  = nullptr;

    do
    {
        stream.avail_in = static_cast<uint32_t>(in_file->read(input_buffer, INFLATE_BUF_SIZE));
        eof             = in_file->eof();
        stream.next_in  = input_buffer;
        do
        {
            stream.avail_out = INFLATE_BUF_SIZE;
            stream.next_out  = output_buffer;

            auto begin = std::chrono::steady_clock::now();
            ret        = isal_inflate(&stream);
            auto end   = std::chrono::steady_clock::now();
            duration += (end - begin);

            out_file->write(output_buffer, INFLATE_BUF_SIZE - stream.avail_out);
        } while (stream.avail_out == 0);
    } while (ret != ISAL_END_INPUT && eof == 0);

    return duration;
}
