#!/bin/bash

set -e

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <module-name> <new-project-folder>"
  echo "Example: $0 github.com/yourusername/myproject myproject"
  exit 1
fi

MODULE_NAME=$1
NEW_PROJECT_FOLDER=$2

echo "Creating new project folder: $NEW_PROJECT_FOLDER"
if [ -d "$NEW_PROJECT_FOLDER" ]; then
  echo "Error: Folder $NEW_PROJECT_FOLDER already exists."
  exit 1
fi

echo "Copying current template project files..."
# Copy all except .git and generate.sh to new folder
rsync -av --progress --exclude '.git' --exclude 'generate.sh' ./ "$NEW_PROJECT_FOLDER"

echo "Replacing placeholders {{.ModuleName}} with $MODULE_NAME ..."
grep -rl '{{.ModuleName}}' "$NEW_PROJECT_FOLDER" | xargs sed -i "s|{{.ModuleName}}|$MODULE_NAME|g"

cd "$NEW_PROJECT_FOLDER"

echo "Running go mod tidy..."
go mod tidy

echo "Done! Your new project is ready at ./$NEW_PROJECT_FOLDER"
