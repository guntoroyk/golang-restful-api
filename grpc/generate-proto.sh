#!/usr/bin/env bash

protoc grpc/proto/category.proto --go_out=plugins=grpc:.
