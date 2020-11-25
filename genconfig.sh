#!/usr/bin/env bash

out=${1:-cmd/ping/kodata/config.js}

echo "writing to $out"

echo "let _URLS = {" > $out
echo '  "global": "https://global.gcping.com",' >> $out
terraform output services | tail -n+2 | sed -e 's/ =/:/g' | sed -e 's/"$/",/g' | sed -e 's/}/};/g' >> $out

cat $out
