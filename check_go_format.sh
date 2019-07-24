#!/bin/bash

git diff-index --quiet HEAD --
if [ $? -ne 0 ]; then
	echo "There's uncommitted changes. This is intended to be run in CI."
	exit 0 # We don't want to have false-positives due to uncommitted changes.
fi

# Run go fmt and make it write the changes into file.
go fmt ./...

git diff-index --quiet HEAD --
if [ $? -eq 0 ]; then
	echo "go fmt test has passed."
	exit 0 # Nothing has been modified.
fi

echo "go fmt test failed. Printing out diff."
git --no-pager diff # Dump out `git diff` so everyone could see what went wrong.
exit 1   # Fail the CI build.
