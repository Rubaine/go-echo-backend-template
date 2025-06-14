#!/bin/sh
set -e
read -p "Module path (ex: github.com/username/project): " module
if [ -z "$module" ]; then
  echo "No module path provided" >&2
  exit 1
fi

go mod edit -module "$module"
find . -name '*.go' -type f -print0 | xargs -0 sed -i "s#example.com/template#$module#g"

go mod tidy

echo "Module initialized to $module"

