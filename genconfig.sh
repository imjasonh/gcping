#!/usr/bin/env bash

out=${1:-cmd/ping/kodata/config.js}

echo "writing to $out"

global=$(terraform output global)

echo "let _URLS = {" > $out
/bin/bash: :wqa: command not found
terraform output services | tail -n+2 | sed -e 's/ =/:/g' | sed -e 's/"$/",/g' | sed -e 's/}/};/g' >> $out

cat $out
