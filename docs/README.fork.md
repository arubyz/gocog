# gocog - generate code for any language, with any language

This is a fork of [natefinch/gocog](https://github.com/natefinch/gocog) with the following changes:

* Added `go.mod` for building with modern versions of Go.

* Added support for [devenv](https://devenv.sh).

* Formatting and linting fixes.

* Added the `--retain` (`-r`) argument to aid debugging by preventing temporary files
  with generator code from being deleted.

* Changed the default generator language to Perl, which is available on virtually
  all platforms.

  * Changed the default value for `--ext` to `pl`.

  * Changed the default value for `--cmd` to `perl`.

  * Changed the default value for `--args` to `%s`.

* Generalized the `--startmark` and `--endmark` arguments and added an additional `--outmark` argument.
  New default values for these arguments enable syntax as follows:
  ```c
  void main()
  {
    // [[[generate]]]
    // $answer = sqrt(42);
    // print "printf(\"the answer is $answer\");";
    // [[[output]]]
    printf("the answer is 6.48074069840786");
    // [[[end]]]
  }
  ```
  The syntax of the original version of `gocog` can be enabled with:
  ```sh
  gocog --startmark '[[[gocog' --outmark 'gocog]]]' --endmark '[[[end]]]'
  ```

* The output of generator code is indented according to the first non-whitespace character
  on the line with the start mark.  This relieves generator code from having to manually
  apply the appropriate indent to each line.

* The logic used to remove indentation and/or comment prefixes from generator code has
  been improved to better support generator code which is column-sensitive (such as perl
  [here documents](https://en.wikipedia.org/wiki/Here_document) delimiters which must
  occur on the first column).  The rules are:

  If the first line of generator code has the same prefix (including indent) as the line
  with the start mark, then all lines of generator code are expected to start with that
  same prefix (which is remove).  This handles generator code in line-oriented comments
  like this:
  ```cpp
  void main()
  {
      // [[[generate]]]
      // print <<done
      // if (true)
      //     printf("Hello, world")
      // done
      // [[[output]]]
      if (true)
          printf("Hello, world")
      // [[[end]]]
  }
  ```

  Otherwise, all lines of generator code are expected to have an indent at least as large
  as the indent of the first line of generator code.  This handles generator code in block
  comments such as this:
  ```c
  void main()
  {
      /* [[[generate]]]
         print <<done
         if (1)
             printf("Hello, world")
         done
         [[[output]]] */
      if (1)
          printf("Hello, world")
      /* [[[end]]] */
  }
  ```
  or this:
  ```c
  void main()
  {
      /* [[[generate]]]
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

You can view the original `README.md` file [here](../README.md).