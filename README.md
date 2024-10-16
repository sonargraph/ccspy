# ccspy
Tool to capture C/C++ compiler command lines for use by [Sonargraph-Architect](https://hello2morrow.com/products/sonargraph/architect9).

## Purpose of the Tool
When Sonargraph analyzes C/C++ code it must know the compiler options used to properly compile a C/C++ compilation unit. The most
important ones are include options (-I) that designate locations where additional header files can be found and macro definitions (-D)
that can influence conditional compilation. While it is possible to figure out these options by manually inspecting compilation logs
it is much easier to collect them automatically via `ccspy`. The tool acts as an intermediary between your build tool (e.g. `make`)
and the compiler. For each compilation unit it records all command line options in a file in the `ccspy` target directory before
redirecting the call to the actual compiler. This target directory should be located parallel to the Sonargraph system directory 
of the system to be analyzed.

## Installation
If you have go installed on your machine you can simply call

`go install github.com/sonargraph/ccspy`

The tool is open source, and you can inspect the code on Github. It is also distributed with your Sonargraph installation 
and can be found in the `bin` directory. On a Mac you have to open the package contents of the Sonargraph app to 
find it. Then make sure that the location of `ccspy` is added to your `PATH` environment variable.

## Integration into your Build
Most C/C++ systems are compiled using tools like `make`,`cmake` or other. To integrate `ccspy` into your build
just configure `ccspy` as your compiler, e.g. in your makefile add a line `CC=ccspy`. Then run a clean build that 
forces all your source files to be compiled via `ccspy`. The result will be one file per compilation unit in
the `ccspy` target directory, which then will be used by `Sonargraph` to retrieve the options needed to for analyzing
each compilation unit correctly.

Keep that configuration as long as you use Sonargrah. If you add new compilation units they will be added automatically
to Sonargraph's analysis.

## Configuration

The tool needs to know which compiler to call and where to store the results (`ccspy` target directory). There 
are 3 modes of configuration:

1. Command line arguments
2. Prefix mode
3. Environment variables

## Command line Arguments
Command line options for `ccspy` always have to be the first parameters in the command line.These are the 
supported command line arguments:
- `-ccspyTargetDir=<path to the target directory>`
- `-ccspyCompiler=<name of your compiler, e.g. gcc>`

Example:

`ccspy -ccspyCompiler=gcc -ccspyTargetDir=/project/ccspy -I../inc -DMACRO=1 test.cpp`

## Prefix Mode
Here you just use `ccspy` before the complete compile command, e.g.:

`ccspy -ccspyTargetDir=/project/ccspy gcc -c -I../inc -DMACRO=1 test.cpp`

In that case you can also define the target directory via environment variable (see below).

## Environment Variables
The tool recognizes the following environment variables:

`CCSPY_CC` default C compiler

`CCSPY_CXX` default C++ compiler

`CCSPY_TARGET_DIR` target directory

## Support
The quickest way to get help is to send an email to `support at hello2morrow.com`.



