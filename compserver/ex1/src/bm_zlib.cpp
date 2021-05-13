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


#include "bm_zlib.h"
#include <zlib.h>

#define BUF_SIZE 262144

std::string bm_zlib::version()
{
    return ZLIB_VERSION;
}

bm::raw_duration bm_zlib::iter_deflate(file_wrapper* in_file, file_wrapper* out_file, int config)
{
    raw_duration duration{};

    int      ret;
    int      flush;
    unsigned have;
    z_stream strm;
    uint8_t  input_buffer[BUF_SIZE];
    uint8_t  output_buffer[BUF_SIZE];

    strm.zalloc = Z_NULL;
    strm.zfree  = Z_NULL;
    strm.opaque = Z_NULL;
    ret         = deflateInit(&strm, config);
    if (ret != Z_OK)
        return raw_duration{0};

    do
    {
        strm.avail_in = static_cast<uInt>(in_file->read(input_buffer, BUF_SIZE));
        if (in_file->error())
        {
            deflateEnd(&strm);
            return raw_duration{0};
        }
        flush        = in_file->eof() ? Z_FINISH : Z_NO_FLUSH;
        strm.next_in = input_buffer;
        do
        {
            strm.avail_out = BUF_SIZE;
            strm.next_out  = output_buffer;

            auto begin = std::chrono::steady_clock::now();
            ret        = deflate(&strm, flush);
            auto end   = std::chrono::steady_clock::now();
            duration += (end - begin);

            if (ret == Z_STREAM_ERROR)
                return raw_duration{0};

            have = BUF_SIZE - strm.avail_out;
            if (out_file->write(output_buffer, have) != have || out_file->error())
            {
                deflateEnd(&strm);
                return raw_duration{0};
            }
        } while (strm.avail_out == 0);
    } while (flush != Z_FINISH);
    deflateEnd(&strm);
    return duration;
}

bm::raw_duration bm_zlib::iter_inflate(file_wrapper* in_file, file_wrapper* out_file)
{
    raw_duration duration{};

    unsigned have;
    z_stream strm;
    uint8_t  input_buffer[BUF_SIZE];
    uint8_t  output_buffer[BUF_SIZE];

    int ret;
    strm.zalloc   = Z_NULL;
    strm.zfree    = Z_NULL;
    strm.opaque   = Z_NULL;
    strm.avail_in = 0;
    strm.next_in  = Z_NULL;
    ret           = inflateInit(&strm);
    if (ret != Z_OK)
        return raw_duration{0};

    do
    {
        strm.avail_in = static_cast<uInt>(in_file->read(input_buffer, BUF_SIZE));
        if (in_file->error())
        {
            inflateEnd(&strm);
            return raw_duration{0};
        }
        if (strm.avail_in == 0)
            break;
        strm.next_in = input_buffer;
        do
        {
            strm.avail_out = BUF_SIZE;
            strm.next_out  = output_buffer;

            auto begin = std::chrono::steady_clock::now();
            ret        = inflate(&strm, Z_NO_FLUSH);
            auto end   = std::chrono::steady_clock::now();
            duration += (end - begin);

            switch (ret)
            {
                case Z_NEED_DICT:
                case Z_DATA_ERROR:
                case Z_MEM_ERROR: (void)inflateEnd(&strm);
                case Z_STREAM_ERROR: return raw_duration{0};
            }

            have = BUF_SIZE - strm.avail_out;
            if (out_file->write(output_buffer, have) != have || out_file->error())
            {
                inflateEnd(&strm);
                return raw_duration{0};
            }
        } while (strm.avail_out == 0);
    } while (ret != Z_STREAM_END);
    inflateEnd(&strm);
    return duration;
}
