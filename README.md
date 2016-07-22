<!--
rm -rf binaries
mkdir -p binaries/
VERSION=0.3.1
for os in darwin linux windows; do
	GOOS=$os GOARCH=$arch go build -o binaries/gargs_${os} main.go
done
-->
gargs
=====

[![Build Status](https://travis-ci.org/brentp/gargs.svg?branch=master)](https://travis-ci.org/brentp/gargs)

Work In Progress:

**gargs** is like **xargs** but it addresses the following limitations in xargs:

+ it keeps the output serialized even when using multiple threads
+ easy to specify multiple arguments with number blocks ({0}, {1}, ...) and {} indicates the entire line.
+ easy to use multiple lines to fill command-template.
+ it defaults to exiting all commands when an error in one of them occurs.

For example, can consume lines of 3 (-n2) and use each line as an `{\d}` command-template filler:

```
$ seq 12 -1 1 | gargs -p 4 -n 3 "sleep {0}; echo {1} {2}"
2 1
5 4
8 7
11 10
```



Install
=======

Download the appropriate binary for your system from [releases](https://github.com/brentp/gargs/releases) into your $PATH.


Limitations
===========

Potentially many processes are writing to the single STDOUT. Currently,
this is handled by reading all output from each process into memory and then
writing that to STDOUT as soon as the previous process has finished writing.
This can lead to large memory use (even if there is no blocking or waiting)
if the output from each command is large.

This can be mitigated by reading the output of each as a stream and sending
to STDOUT in chunks of lines but this assumes all output is line-wise.

We may also be able to mitigate by using io.PipeReader.

Example
=======
Let's say we have a file `t.txt` like:
```
chr1	22 33
chr2 22 33
chr3 22	33
chr4	22	33
```
That has a mixture of tabs and spaces. We can convert each line to chrom:start-end format with:

```
$ cat t.txt | gargs --sep "\s+" -p 2 "echo '{0}:{1}-{2}' full-line: \'{}\'"
chr2:22-33 full-line: 'chr2 22 33'
chr1:22-33 full-line: 'chr1 22 33'
chr3:22-33 full-line: 'chr3 22 33'
chr4:22-33 full-line: 'chr4 22 33'
```

In this case, we're using **2** processes to run this in parallel which will make more of a difference
if we do something time-consuming rather than `echo`. The output will be kept in the order dictated by
`t.txt` even if the processes finish in a different order. This is sometimes at the expense of parallelization
efficiency.


Usage
=====

```
usage: gargs_race [--procs PROCS] [--nlines NLINES] [--sep SEP] [--verbose] [--continue-on-error] [--ordered] [--dry-run] COMMAND

positional arguments:
  command                command to execute

options:
  --procs PROCS, -p PROCS
                         number of processes to use [default: 1]
  --nlines NLINES, -n NLINES
                         number of lines to consume for each command. -s and -n are mutually exclusive. [default: 1]
  --sep SEP, -s SEP      regular expression split line with to fill multiple template spots default is not to split. -s and -n are mutually exclusive.
  --verbose, -v          print commands to stderr before they are executed.
  --continue-on-error, -c
                         report errors but don't stop the entire execution (which is the default).
  --ordered, -o          keep output in order of input; default is to output in order of return which greatly improves parallelization.
  --dry-run, -d          print (but do not run) the commands
  --help, -h             display this help and exit
```

TODO
====

+ [X] final exit code is the largest of any seen exit code even with -c
+ [X] dry-run
+ [ ] combinations of `-n` and `--sep`.
+ [ ] for example, we are sending regions to bcftools view or tabix. It's faster to send multiple
      queries to each rather than starting a new process for each one.
      if we do 'bcftools view {} {} {}' then we'll get an error at the last round if the input is
      not divisible by 3. For this case, we should be able to just issue a warning.

	  Actually, can currently do: `seq 10 | ./gargs -n 3 "echo {0} {1} {2}"`

