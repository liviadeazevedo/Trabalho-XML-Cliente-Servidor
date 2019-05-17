#!/bin/bash

cd serverLog;
go build && go install;

cd ../serverLogic;
go build && go install;

cd ../serverConnection;
go build && go install;