# gocog - generate code for any language, with any language

This is a fork of [natefinch/gocog](https://github.com/natefinch/gocog) with the following changes:

* Added `go.mod` for building with modern versions of Go.

* Added support for [devenv](https://devenv.sh).

* Formatting and linting fixes.

* Removed the default value for `--cmd`, requiring that the scripting language always be
  explicitly specified.

* Changed the default value for `--args` from `run '%s'`, which is specific to Go, to just `'%s'`.

* Generalized the `--startmark` and `--endmark` arguments and added an additional `--outmark` argument.
  New default values for these arguments enable syntax as follows:
  ```c
  void main()
  {
    // [[[generate]]]
    // const answer = 42;
    // console.log(`  printf("the answer is ${answer}");`);
    // [[[output]]]
    printf("the answer is 42");
    // [[[end]]]
  }
  ```
  The syntax of the original version of `gocog` can be enabled with:
  ```sh
  gocog --startmark '[[[gocog' --outmark 'gocog]]]' --endmark '[[[end]]]'
  ```

You can view the original `README.md` file [here](../README.md).