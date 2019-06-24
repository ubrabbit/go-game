#!/bin/bash

go test -bench=. -benchmem -benchtime="1s"
