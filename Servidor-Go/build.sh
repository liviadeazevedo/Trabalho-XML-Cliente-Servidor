#!/bin/bash
(
    cd serverLogic &&
    go build &&
    go install
) &&
(
    cd serverConnection
    go build &&
    go install
)