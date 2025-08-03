#!/bin/bash

# Define the possible script names
scripts=("build" "build-ts")

# Check each possible script name
for script in "${scripts[@]}"; do
  if npm run | grep -q " $script"; then
    echo "Running npm script: $script"
    npm run "$script"
    exit 0
  fi
done

# If no matching script is found
echo "Error: No suitable build script found."
exit 1