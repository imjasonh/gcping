#!/usr/bin/env bash

out=${1:-cmd/ping/kodata/config.js}

echo "writing to $out"

echo -n "let _URLS = " > $out
terraform output services | sed -e 's/ =/:/g' | sed -e 's/"$/",/g' | sed -e 's/}/};/g' >> $out
