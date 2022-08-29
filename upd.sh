#!/usr/bin/bash

for i in pkg/*; do
  cd $i && go get -u && cd ../..
done