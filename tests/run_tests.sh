#!/bin/sh
set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Running hurl tests against $BASE_URL"
echo "========================================"

for file in "$SCRIPT_DIR"/health.hurl "$SCRIPT_DIR"/categories.hurl "$SCRIPT_DIR"/products.hurl; do
    if [ ! -f "$file" ]; then
        continue
    fi
    echo ""
    echo "--- $file ---"
    hurl --test "$file" --variable base_url="$BASE_URL"
done

echo ""
echo "All tests passed."
