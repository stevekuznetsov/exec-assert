# exec-assert

`exec-assert` is a testing framework for bash written in Go

## Installation

`exec-assert` can be built from the repository root with:

```sh
$ make build
```

## Usage

`exec-assert` consumes a bash command and executes it, collecting the output to `stdout` and `stderr`, allowing the program to make assertions about the result of the command execution as well as the content of its output. 

`exec-assert` can run the bash command once and assert that:
 * the command succeeds
 * the command fails
 * the command's output contains text or a regular expression
 * the command's output does not contain text or a regular expression
 * any combination of a result assertion and any number of output assertions

`exec-assert` will succeed only if all of the assertions are met. 

`exec-assert` can furthermore run the bash command with some regular interval, until the result of execution fulfills all of the assertions or a timeout.

### Examples

To test that a command (`date`) executes successfully:
```sh
$ exec-assert 'date'
executing "date" once, expecting success 
SUCCESS after 0.002s: executing "date" once, expecting success 
```

To test that a command (`date`) executes successfully and its output contains a phase (`Wed`):
```sh
$ exec-assert --output contains --test 'Wed' 'date'
executing "date" once, expecting success and output that contains "Wed"
SUCCESS after 0.002s: executing "date" once, expecting success and output that contains "Wed"
```

To test that a command (`date`) executes successfully and its output contains a phase (`Wed`) but doesn't contain a regular expression (`201[0-4]`):
```sh
$ exec-assert --output 'contains,excludes' --test 'Wed,201[0-4]' --delimiter ',' 'date'
executing "date" once, expecting success and output that contains "Wed", and doesn't contain "201[0-4]"
SUCCESS after 0.004s: executing "date" once, expecting success and output that contains "Wed", and doesn't contain "201[0-4]"
```

To test that a command (`grep`) fails to execute:
```sh
$ exec-assert --result failure 'grep'
executing "grep" once, expecting failure 
SUCCESS after 0.003s: executing "grep" once, expecting failure
```

To run a command (`date`) without regard to its return code until its output contains a regular expression (`\:2{2}`), choosing verbose execution to see the command's output:
```sh
$ $ ./exec-assert --execute until --result ambivalent --output contains --test '\:2{2}' -v 'date'
executing "date" every 0.200s for 60.000s, or until success and output that contains `\:2{2}`
SUCCESS after 1.611s: executing "date" every 0.200s for 60.000s, or until success and output that contains `\:2{2}`
Command output to stdout: 
5x  Wed Dec 30 15:49:20 MST 2015
  --
3x  Wed Dec 30 15:49:21 MST 2015
  --
1x  Wed Dec 30 15:49:22 MST 2015
Command did not output to stderr.
```

More examples of usage can be found in the [integration test](test/cmd.sh).

### Correctly Quoting Text and Variables
To run a command that doesn't contain any quoted text, quote the argument to `exec-assert` with single or double quotes:
```sh
$ exec-assert 'ls -lAhR'
$ exec-assert "ls -lAhR"
```

To run a command that contains literal text, quote the arguemnt to `exec-assert` with double quotes and the literal text with single quotes:
```sh
$ exec-assert "echo 'string'"
```

To run a command that contains a bash variable, the argument to `exec-assert` *must* be quoted with double quotes for the variable to be expanded correctly:
```sh
$ myvar=value
$ exec-assert "echo 'expression containing ${myvar}'"
```

To run a command that contains something that looks like a bash variable, but isn't, or a bash variable that you do not want to be expanded, escape the `$` with a forward slash in the argument to `exec-assert`:
```sh
$ myvar=value
$ exec-assert "echo '\$myvar=${myvar}'"
```

### Contributing

Contributions are welcome to this repository. All pull requests will be tested by applying `go vet` for linting, `go test` for unit testing, and running `test/cmd.sh` for integation testing. All of these tests can be run locally with:
```sh
$ make verify
```