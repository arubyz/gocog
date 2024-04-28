# gocog - generate code for any language, with any language

This is a fork of [natefinch/gocog](https://github.com/natefinch/gocog) with the following changes:

* Added `go.mod` for building with modern versions of Go.

* Added support for [devenv](https://devenv.sh).

* Formatting and linting fixes.

* Added the `--retain` (`-r`) option to aid debugging by preventing temporary files
  with generator code from being deleted.

* When operating over multiple files using `@filelist.txt` syntax, the files listed
  are interpreted relative to the directory containing `filelist.txt`, not the current
  directory when `gocog` is run.

* Removed the default value of `go` for `--cmd`, requiring it to be explicitly specified 
  either on the command line or via the start mark (see below).  Also changed the default
  for `--ext` to `.txt` and for `--args` to `%s`, which are defaults that work well for
  a broad set of (non-Go) languages processors.

* Generalized the `--startmark` and `--endmark` arguments and added an additional `--outmark`
  argument.  The specification of markers was also changed from literal strings to regular
  expressions.  For `--startmark` specifically, if the regular expression defines a group
  then the value that group matches is used as the command to process the generator code
  for that block.  These changes enable syntax as follows:
  ```c
  void main()
  {
    // [[[generate perl]]]
    // $answer = sqrt(42);
    // print "printf(\"the answer is $answer\");";
    // [[[output]]]
    printf("the answer is 6.48074069840786");
    // [[[end]]]
  }
  ```
  The syntax of the original version of `gocog` can be enabled with:
  ```sh
  gocog --startmark '\[\[\[gocog' --outmark 'gocog\]\]\]' --endmark '\[\[\[end\]\]\]'
  ```

* The output of generator code is indented according to the first non-whitespace character
  on the line with the end mark.  This relieves generator code from having to manually
  apply the appropriate indent to each line.  Using the line with the start mark or 
  output mark instead doesn't always work, since these lines may be inside a multi-line
  block comment with additional indentation.  The end mark though should always be in a
  single-line comment, regardless of whether it's a block or line comment, so that line's
  indent should always be indicative of the desired indent for the generated lines.

* The logic used to remove indentation and/or comment prefixes from generator code has
  been improved to better support generator code which is column-sensitive (such as perl
  [here document](https://en.wikipedia.org/wiki/Here_document) delimiters which must
  occur on the first column).  The rules are:

  1. If the first line of generator code has the same prefix (including indent) as the line
     with the start mark, then all lines of generator code are expected to start with that
     same prefix (which is remove).  This handles generator code in line-oriented comments
     like this:
     ```cpp
     void main()
     {
         // [[[generate perl]]]
         // print <<done
         // if (true)
         //     printf("Hello, world");
         // done
         // [[[output]]]
         if (true)
             printf("Hello, world");
         // [[[end]]]
     }
     ```

  2. Otherwise all lines of generator code are expected to have an indent at least as large
     as the indent of the first line of generator code, and each line is un-indented according
     to the indent of the first line of generator code.  This handles generator code in block
     comments such as this:
     ```c
     void main()
     {
         /* [[[generate perl]]]
            print <<done
            if (1)
                printf("Hello, world");
            done
            [[[output]]] */
         if (1)
             printf("Hello, world");
         /* [[[end]]] */
     }
     ```
     or this:
     ```c
     void main()
     {
         /* [[[generate perl]]]
         print <<done
         if (1)
             printf("Hello, world");
         done
         [[[output]]] */
         if (1)
             printf("Hello, world");
         /* [[[end]]] */
     }
     ```

* Added the `--extraline` (`-L`) option to cause the line after the line with the output
  mark to also be considered part of the output mark.  This accommodates generator code
  in block comments where the comment end delimiter is on its own line, like this:
  ```c
  void main()
  {
      /*
      [[[generate perl]]]
      print <<done
      if (1)
          printf("Hello, world");
      done
      [[[output]]]
      */
      if (1)
          printf("Hello, world");
      /* [[[end]]] */
  }
  ```

You can view the original `README.md` file [here](../README.md).