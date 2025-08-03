#!/bin/bash

# Define the possible script names
scripts=("dev" "start:dev" "start-dev")

# Check each possible script name
for script in "${scripts[@]}"; do
  if npm run | grep -q " $script"; then
    echo "Running npm script: $script"
    npm run "$script"
    exit 0
  fi
done

# If no matching script is found
echo "Error: No suitable dev script found."
exit 1