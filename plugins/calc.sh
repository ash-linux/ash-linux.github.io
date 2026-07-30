#!/bin/bash
QUERY="${1#calc:}"
echo "$QUERY" | bc -l
