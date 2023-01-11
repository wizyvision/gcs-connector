#!/usr/bin/env bash

# Exports a flat yaml key:value pairs as environment variables
# https://gist.github.com/bermi/14f6584215dd1524754077b839c1a775
#
# To source a yaml file as environment variables run
#
#     . yaml-source env.yml
#
eval $(
  cat $1 |
  # limit only to lines that contain a valid KEY
  grep -E '^[A-Za-z0-9_ ]+:' |
  # replace `key:'values with quotes'` with `export key='value with quotes'`
  sed -E 's/^([A-Za-z0-9_]+) *: *([\x27"])/export \1=\2/' |
  # quote items that missed the quotes
  sed -E 's/^([A-Za-z0-9_]+) *: *(.+)/export \1="\2"/'
)