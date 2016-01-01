#!/bin/bash

# This file contains a series of CLI tests for exec-assert

set -o errexit
set -o nounset
set -o pipefail

function exit_report() {
	out=$?
	if [[ ${out} -ne 0 ]]; then
		echo "[FAIL] Tests failed!"
	else 
		echo "[SUCCESS] Tests succeeded!"
	fi
}

trap exit_report EXIT INT TERM

echo "[INFO] Building 'exec-assert' binary..."
go build .
echo "[INFO] Built 'exec-assert' binary successfuly"

echo "[INFO] Beginning tests..."

# Simple tests of successful invocations
./exec-assert 'pwd'
./exec-assert --result success 'pwd'
./exec-assert --result success --output contains --test '/' 'pwd'
./exec-assert --result success --output excludes --test '#' 'pwd'
./exec-assert --output 'contains,excludes' --test '/,#' --delimiter ',' 'pwd'
./exec-assert --result failure 'grep'
./exec-assert --result failure --output contains --test 'for more information' 'grep'
./exec-assert --result failure --output excludes --test 'bogus text' 'grep'
./exec-assert --result failure --output 'contains,excludes' --test 'for more information#bogus text' --delimiter '#' 'grep'
./exec-assert --result ambivalent 'exit 0'
./exec-assert --result ambivalent 'exit 1'

# Simple tests of failing invocations
if ./exec-assert --result failure 'date'; then
	exit 1
fi
if ./exec-assert --result success --output excludes --test '[0-9]' 'date'; then
	exit 1
fi
if ./exec-assert --output contains --test '1999' 'date'; then
	exit 1
fi

# Exection strategy "until" tests
./exec-assert --execute until --output contains --test ':[0-9]5 ' --timeout 11s --interval 1s -v 'date' # so that only seconds can fulfill and re-tries happen 
./exec-assert --execute until --result success 'pwd'

if ./exec-assert --execute until --result success --timeout 2s 'exit 1'; then
	exit 1
fi
if ./exec-assert --execute until --output excludes --test 'hello' --timeout 2s 'echo hello'; then
	exit 1
fi

# Complex command tests
# Pipes
./exec-assert 'echo "hello" | grep "hello"'

./exec-assert --output contains --test 'exec-assert\n' 'echo "-1" | xargs ls'

# Variable expansion
test_var=TEST
./exec-assert --output excludes --test 'TEST' 'echo "${test_var}"' # we aren't running in a subshell so we can't expand $test_var
./exec-assert --output contains --test 'TEST' "echo ${test_var}" # if we get the expanded var in the first place, everything's fine
unset test_var

# Multiple statements
./exec-assert --result success --output contains --test 'hello' "echo 'hello'; exit 0"

# Curly-brace expansions
./exec-assert --output 'contains,contains' --test 'options,interfaces' --delimiter ',' -v 'ls pkg/cmd/{option,interface}s.go'

# Logic
./exec-assert 'if true; then exit 0; else exit 1; fi'

# Integer arithmetic
./exec-assert --output contains --test '19' 'echo $(( 20 - 1 ))'

# Stdout/Stderr redirects
./exec-assert --output excludes --test '[A-Za-z0-9/]' 'pwd 1>/dev/null'
./exec-assert --result failure --output excludes --test 'for more information' 'grep 2>/dev/null'

# `Here` documents and strings
./exec-assert 'grep hello <<EOF
hello
EOF
'

./exec-assert 'grep hello <<< hello'
