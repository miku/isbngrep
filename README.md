README
======

A simple command line tool to extract normalized ISBNs from arbitrary text.

Help
----

    $ ./isbngrep -h
    NAME:
       isbngrep - find ISBNs in texts

    USAGE:
       isbngrep [global options] command [command options] [arguments...]

    VERSION:
       1.0.0

    COMMANDS:
       help, h  Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --verbose        be verbose
       --uniq, -u       return a uniq list
       --only-10, -0    only ISBN10
       --only-13, -3    only ISBN13
       --version, -v    print the version
       --help, -h       show help


Build
-----

    $ git clone git@github.com:miku/isbngrep.git
    $ cd isbngrep
    $ go get github.com/codegangsta/cli
    $ make
    $ cat test.txt | ./isbngrep --uniq
    0486653552
    0486281523
    0486234002
    0486270785
    ...
