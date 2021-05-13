## Prerequisite

Install the dependencies:
* a c++14 compliant c++ compiler
* cmake >= 3.1
* git
* autogen
* autoconf
* automake
* yasm and/or nasm
* libtool
* boost's "Filesystem" library and headers
* boost's "Program Options" library and headers
* boost's "String Algo" headers

```
sudo apt-get update
sudo apt-get install gcc g++ make cmake git zlib1g-dev autogen autoconf automake yasm nasm libtool libboost-all-dev 
```

Also needed are the latest versions of isa-l and zlib. The `get_libs.bash` script can be used to get them.
The script will download the two libraries from their official GitHub repository, build them, and install them in `./libs/usr`.

* `bash ./libs/get_libs.bash`

## Build

* `mkdir <build-dir>`
* `cd <build-dir>`
* `cmake -DCMAKE_BUILD_TYPE=Release $OLDPWD`
* `make`

## Run

* `cd <build-dir>`
* `./ex1 --help`

## Usage

```
Usage: ./ex1 [--help] [--folder <path>]... [--file <path>]... :
  --help                    display this message
  --file path               use the file at 'path'
  --folder path             use all the files in 'path'
  --zlib-levels n,...       coma-separated list of compression level [1-9]
  --semidyn-config flag,... coma-separated list of flags for the semi-dynamic
                            compression ('stateful', 'stateless') [stateful]
```

* `--file` and `--folder` can be used multiple times to add more files to the benchmark
* `--folder` will look for files recursively
* the default `--zlib-level` is `6`

Test corpuses are available online (e.g.: [Calgary](http://corpus.canterbury.ac.nz/descriptions/#calgary), [Silesia](http://sun.aei.polsl.pl/~sdeor/index.php?page=silesia)).
The `--folder` option can be use to easily benchmark them: `./ex1 --folder /path/to/corpus/folder`.


## Example run

```
$> ./ex1 --file test/file1 --file test/file2
[Info   ] Using isa-l      2.17.0
[Info   ] Using zlib       1.2.11
[Info   ] Running benchmarks at 2902.71 MHz
[Warning] CPU scaling is enabled. Real time measurements will be noisy
[Info   ] Benchmarks starting...
[Info   ] ( 1/12) Compressing test/file1 (10.2 MiB) with isa-l (static)...
[Info   ] ( 2/12) Decompressing ...
[Info   ] ( 3/12) Compressing test/file1 (10.2 MiB) with isa-l (semi-dyn)...
[Info   ] ( 4/12) Decompressing ...
[Info   ] ( 5/12) Compressing test/file1 (10.2 MiB) with zlib...
[Info   ] ( 6/12) Decompressing ...
[Info   ] ( 7/12) Compressing test/file2 (5.3 MiB) with isa-l (static)...
[Info   ] ( 8/12) Decompressing ...
[Info   ] ( 9/12) Compressing test/file2 (5.3 MiB) with isa-l (semi-dyn)...
[Info   ] (10/12) Decompressing ...
[Info   ] (11/12) Compressing test/file2 (5.3 MiB) with zlib...
[Info   ] (12/12) Decompressing ...

test/file2 (5.3 MiB)
---------------------------------------------------------------------------------------------------------------------------------------------
 Library         | Config    | Compression                                                      | Decompression
                 |           | Ratio              | Real Time            | CPU Time             | Real Time            | CPU Time
-----------------|-----------|--------------------|----------------------|----------------------|----------------------|----------------------
isa-l (static)   |         - |   19.84 % (x 1.00) |    7867 us (x  1.00) |    9377 us (x  1.00) |    5558 us (x  1.00) |    8873 us (x  1.00)
isa-l (semi-dyn) |  stateful |   19.66 % (x 0.99) |    9310 us (x  1.18) |   11080 us (x  1.18) |    5964 us (x  1.07) |    9185 us (x  1.04)
zlib             |   level 6 |   12.87 % (x 0.65) |  121779 us (x 15.48) |  123437 us (x 13.16) |   11677 us (x  2.10) |   15032 us (x  1.69)

test/file1 (10.2 MiB)
---------------------------------------------------------------------------------------------------------------------------------------------
 Library         | Config    | Compression                                                      | Decompression
                 |           | Ratio              | Real Time            | CPU Time             | Real Time            | CPU Time
-----------------|-----------|--------------------|----------------------|----------------------|----------------------|----------------------
isa-l (static)   |         - |   46.14 % (x 1.00) |   34410 us (x  1.00) |   39557 us (x  1.00) |   26624 us (x  1.00) |   33825 us (x  1.00)
isa-l (semi-dyn) |  stateful |   44.67 % (x 0.97) |   37394 us (x  1.09) |   42403 us (x  1.07) |   27565 us (x  1.04) |   34847 us (x  1.03)
zlib             |   level 6 |   37.99 % (x 0.82) |  774417 us (x 22.51) |  780019 us (x 19.72) |   47542 us (x  1.79) |   54768 us (x  1.62)

Average results:
---------------------------------------------------------------------------------------------------------------------------------------------
 Library         | Config    | Compression                                                      | Decompression
                 |           | Ratio              | Real Time            | CPU Time             | Real Time            | CPU Time
-----------------|-----------|--------------------|----------------------|----------------------|----------------------|----------------------
isa-l (static)   |         - | x             1.00 | x               1.00 | x               1.00 | x               1.00 | x               1.00
isa-l (semi-dyn) |  stateful | x             0.98 | x               1.14 | x               1.13 | x               1.05 | x               1.03
zlib             |   level 6 | x             0.74 | x              18.99 | x              16.44 | x               1.94 | x               1.66
```


## Interpreting the output

#### CPU scaling warning

`[Warning] CPU scaling is enabled. Real time measurements will be noisy`

Most modern CPUs have some sort of CPU scaling. This result in their frequencies changing dynamically.
In order to get the best result from this benchmark, CPU scaling should be disabled.

#### Individual file results

For each file, for each library, the output contains:
* the library name
* any specific configuration for that library
* the compression ratio: the value is computed such as 'lower is better'. `75%` means that the compressed file is 75% the size of the original.
* the real-time and cpu-time for the compression, displayed in microseconds ('us') or nanoseconds ('ns').
* the real-time and cpu-time for the decompression. The decompression is performed on the file compressed with that same library.

