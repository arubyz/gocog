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

* When running a command to process generator code, the current directory is set to
  the same directory as the input file.  This makes it easy for generator code to
  reference other files which are in the same directory as the input file.
  
* Removed the `--cmd` argument, requiring the value to be explicitly specified via the
  `--genstart` regexp (see below).  Also changed the default for `--ext` to `.txt` and for `--args`
  to `%s`, which are defaults that work well for a broad set of (non-Go) languages processors.

* Generalized the `--startmark` and `--endmark` arguments by replacing them with four different
  marker regexp definitions:

  * `--genstart`: The first line of generator code is read immediately after the line containing this regexp.
    If this regexp defines a group then the value that group matches is used as the command to process
    the generator code for that block.

  * `--genend`: The last line of generator code is read immediately before the line containing this regexp.

  * `--outstart` (optional): If defined, the first line of output will be inserted immediately after
    the line containing this regexp.  If not defined, the first line of output will be inserted immediately
    after the line containing the `--genend` regexp.

  * `--outend`: The last line of output will be inserted immediately before the line containing this regexp.

  These new arguments and their default values enable syntax as follows:

  ```cpp
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

  And setting `--outstart` to something like `\*\/` supports multi-line comment delimiters
  being on their own line as follows:

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

  These arguments also allow the generator code and output to be separated, eg:

  ```cpp
  // [[[generate perl]]]
  // print <<done
  // if (1)
  //     printf("Hello, world");
  // done
  // [[[end]]]
  void main()
  {
      // [[[output]]]
      if (1)
          printf("Hello, world");
      // [[[end]]]
  }
  ```

  The syntax of the original version of `gocog` can be enabled with:

  ```sh
  gocog --genstart '\[\[\[gocog' --genend 'gocog\]\]\]' --outend '\[\[\[end\]\]\]'
  ```

* The output of generator code is indented according to the first non-whitespace character
  on the line with the `--outend` regexp.  This relieves generator code from having to manually
  apply the appropriate indent to each line.  Using instead the line with the `--genstart` or 
  `--genend` regexp doesn't always work, since these lines may be inside a multi-line
  block comment with additional indentation.  The `--outend` regexp though should always be in a
  single-line comment, regardless of whether it's a block or line comment, so that line's
  indent should always be indicative of the desired indent for the generated lines.

* The logic used to remove indentation and/or comment prefixes from generator code has
  been improved to better support generator code which is column-sensitive (such as perl
  [here document](https://en.wikipedia.org/wiki/Here_document) delimiters which must
  occur on the first column).  The rules are:

  1. If the first line of generator code has the same prefix (including indent) as the line
     with the `--genstart` regexp, then all lines of generator code are expected to start with that
     same prefix (which is removed).  This handles generator code in line-oriented comments
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

* The `--ext` argument has been replaced by `--genfile` which provides a more generic
  way to specify the the full path of the temporary file containing generator code.
  Within this template, ths follow variables are expanded:

  * `$TMP`: The absolute path to the OS-defined temporary file directory

  * `$DIR`: The directory containing the input file

  * `$FILE`: The file name and extension (without directory) of the input file

  * `$CTR`: A counter value which is incremented for each generator block in a given
    input file

* A new `--outfile` argument has been added which specifies a template for the full path
  to the temporary output file.  Within this template, ths follow variables are expanded:

  * `$TMP`: The absolute path to the OS-defined temporary file directory

  * `$DIR`: The directory containing the input file

  * `$FILE`: The file name and extension (without directory) of the input file

You can view the original `README.md` file [here](../README.md).