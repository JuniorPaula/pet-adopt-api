#!/bin/bash

# verifica se a pasta bin existe
if [ ! -d "./bin" ]; then
  mkdir bin
fi

# Run the application
go build -o ./bin/app ./cmd/api/*.go
./bin/app