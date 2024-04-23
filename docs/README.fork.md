# gocog - generate code for any language, with any language

This is a fork of [natefinch/gocog](https://github.com/natefinch/gocog) with the following changes:

* Added `go.mod` for building with modern versions of Go.

* Added support for [devenv](https://devenv.sh).

* Formatting and linting fixes.

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

* Generator code is un-indented according to the column at which the start mark is found,
  rather than the column of the first non-whitespace character on the line with the start
  mark. This enables generator code which is column-sensitive (eg,
  [here documents](https://en.wikipedia.org/wiki/Here_document)) to be processed correctly.

* The output of generator code is indented according to the first non-whitespace character
  on the line with the start mark.  This relieves generator code from having to manually
  apply the appropriate indent to each line.

You can view the original `README.md` file [here](../README.md).