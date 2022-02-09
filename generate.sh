#!/bin/bash

protoc loggerpb/logger.proto --go_out=plugins=grpc:.
