# ccspy
Tool to capture C/C++ compiler command lines for use by [Sonargraph-Architect](https://hello2morrow.com/products/sonargraph/architect9).

## Purpose of the Tool
When Sonargraph analyzes C/C++ code it must know the compiler options used to properly compile a C/C++ compilation unit. The most
imortant ones are include options (-I) that designate locations where additional header files can be found and macro defiitions (-D)
that can influence conditional compilation. While it us possible to figure out these options by manually inspecting compilation logs
it is much easier to collect automatically them via `ccspy`. The tool acts as an intermediary between your build tool (e.g. `make`)
and the compiler. For each compilation unit it records all command line options in a file in the `ccspy` target directory before
redirecting the call to the actual compiler.  This target directory should be located parallel to the Sonargraph system directory 
of the system to be analyzed.
