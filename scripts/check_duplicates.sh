#!/bin/bash

# Find all type declarations and check for duplicates
echo "Checking for duplicate type declarations..."

find pkg/models -name "*.go" -exec grep -l "type.*struct" {} \; | while read file; do
    types=$(grep "^type.*struct" "$file" | awk '{print $2}')
    for type in $types; do
        count=$(find pkg/models -name "*.go" -exec grep -l "^type $type struct" {} \; | wc -l)
        if [ $count -gt 1 ]; then
            echo "Duplicate type $type found in:"
            find pkg/models -name "*.go" -exec grep -l "^type $type struct" {} \;
        fi
    done
done