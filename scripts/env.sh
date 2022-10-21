#!/bin/bash

nprocs=$(getconf _NPROCESSORS_ONLN)

# GO
GOEXEC=${GOEXEC:-"go"}
GOROOT=$GOROOT

REPORT=${REPORT:-"$(date +%F-%H-%M)"}
tee_cmd="tee -a output/${REPORT}.log"
