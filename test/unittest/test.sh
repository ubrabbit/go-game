#!/bin/bash

ulimit -n 65535

go test -cover -v -bench=.
