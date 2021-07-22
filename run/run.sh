#!/bin/bash

cd ../ && go build -o chat main.go

./chat -name=wxy -org=china && ./chat -name=fsj -org=china && ./chat -name=son -org=china && ./chat -name=daughter -org=china

