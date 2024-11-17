#!/bin/bash

# Define the possible script names
scripts=("test:integration" "test")

# Store all additional arguments
args="${@}"

# Check each possible script name
for script in "${scripts[@]}"; do
  if npm run | grep -q " $script"; then
    echo "Running npm script: $script ${args}"
    npm run "$script" -- ${args}
    exit 0
  fi
done

# If no matching script is found
echo "Error: No suitable test:integration script found."
exit 1