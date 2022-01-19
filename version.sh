#!/bin/bash

# Version argument, should be something like '0.0.1'. Will be added to '.version' file.
VER="$1"
# Value to use as the git tag. Will end up being somethig like 'v0.0.1'.
TAG=$VER
# Version file full path.
FILE="$(pwd)/.version"

CURRENT_TAG=$(git describe --abbrev=0 --tags 2>/dev/null)

DIFF="$(git --no-pager diff)"

[ ! -f $FILE ] && touch $FILE

CURRENT="$(cat $FILE)"

# Ensure there are no unstaged changes.
if [[ $DIFF != "" ]]; then
    echo "There are unstaged changes. Commit or stash any changes and try again."
    exit 1
fi

# Ensure a version is provided before continuing.
if [[ $VER == "" ]]; then
    echo "Provide a version" 1>&2
    exit 1
fi

# Remove prepended 'v' from version argument, if it was provided.
if [[ "$VER" == *"v"* ]]; then
    ver=$(echo "$VER" | cut -d 'v' -f 2)
    VER="$ver"
fi

# Append 'v' to tag variable, if 'v' was not already prepended.
if [[ "$TAG" != *"v"* ]]; then
    TAG="v$TAG"
fi

if [[ "$CURRENT" == "$VER" ]]; then
    echo "Version is already $VER"
    exit 1
elif [[ "$CURRENT_TAG" == "$TAG" ]]; then
    git_del_tag="$(git tag -d $TAG 1>/dev/null 2>&1)"

    if [[ "$git_del_tag" != "" ]]; then
        echo -e "Error deleting git tag $TAG:\n$git_del_tag"
        exit 1
    fi
fi

echo $VER >$FILE
echo "Added version $VER to $FILE"

GIT_ADD_ERR="$(git add $FILE 1>/dev/null 2>&1)"

# Print any errors from git and exit.
if [[ "$GIT_ADD_ERR" != "" ]]; then
    echo -e "git add error:\n$GIT_ADD_ERR" 1>&2
    exit 1
fi

GIT_COMMIT_ERR="$(git commit -m "Release $TAG" 1>/dev/null 2>&1)"

# Print any errors from git and exit.
if [[ "$GIT_COMMIT_ERR" != "" ]]; then
    echo -e "git commit error:\n$GIT_COMMIT_ERR" 1>&2
    exit 1
fi

# Capture stderr from adding git tag.
GIT_TAG_ERR="$(git tag $TAG 2>&1)"

# Print any errors from git and exit.
if [[ "$GIT_TAG_ERR" != "" ]]; then
    echo -e "git tag error:\n$GIT_TAG_ERR" 1>&2
    exit 1
else
    echo "Added git tag $TAG"
fi

exit 0
