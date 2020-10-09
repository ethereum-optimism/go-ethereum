#!/bin/sh

cd /go-ethereum

echo "*** Running Linter ***"
go run build/ci.go lint || { echo 'linter failed' ; exit 1; }

echo "*** Linter succeeded ***"

echo "*** Running Tests ***"
go run build/ci.go test -coverage $TEST_PACKAGES || { echo 'tests failed' ; exit 1; }

echo "*** Tests Passed ***"
