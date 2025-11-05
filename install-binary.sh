#!/usr/bin/env bash

while [ $# -gt 0 ]; do
  case "$1" in
    --version*|-v*)
      if [[ "$1" != *=* ]]; then shift; fi
      VERSION="${1#*=}"
      ;;
    *)
      >&2 printf "Error: Invalid argument\n"
      exit 1
      ;;
  esac
  shift
done
if [ -z $VERSION ]; then
    VERSION='latest'
fi

PROJECT_NAME="helm-ssm"
PROJECT_GH="tutti-ch/$PROJECT_NAME"
eval $(helm env)

if [[ $SKIP_BIN_INSTALL == "1" ]]; then
  echo "Skipping binary install"
  exit
fi

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
    armv5*) ARCH="armv5";;
    armv6*) ARCH="armv6";;
    armv7*) ARCH="armv7";;
    aarch64) ARCH="arm64";;
    x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(echo $(uname)|tr '[:upper:]' '[:lower:]')
  case "$OS" in
    # Msys support
    msys*) OS='windows';;
    # Minimalist GNU for Windows
    mingw*) OS='windows';;
  esac
}

# verifySupported checks that the os/arch combination is supported for
# binary builds.
verifySupported() {
    local supported="linux-amd64\ndarwin-amd64\ndarwin-arm64\nwindows-amd64\nlinux-386\nlinux-arm64\nwindows-386\nwindows-arm64"
  if ! echo "${supported}" | grep -q "${OS}-${ARCH}"; then
    echo "No prebuild binary for ${OS}-${ARCH}."
    exit 1
  fi

  if ! type "curl" > /dev/null && ! type "wget" > /dev/null; then
    echo "Either curl or wget is required"
    exit 1
  fi
}

# getDownloadURL checks the latest available version.
getDownloadURL() {
  # Use the GitHub API to find the latest version for this project.
  if [ $VERSION = 'latest' ]; then
    local latest_url="https://api.github.com/repos/$PROJECT_GH/releases/$VERSION"
    local response=""

    # Try authenticated request first if GITHUB_TOKEN is available
    if [ -n "$GITHUB_TOKEN" ]; then
      echo "Using authenticated GitHub API request"
      if type "curl" > /dev/null; then
        response=$(curl -sL -H "Authorization: Bearer $GITHUB_TOKEN" "$latest_url")
      elif type "wget" > /dev/null; then
        response=$(wget -q --header="Authorization: Bearer $GITHUB_TOKEN" -O - "$latest_url")
      fi

      # Check if authentication failed (response contains "Bad credentials" or other error)
      if echo "$response" | grep -q "Bad credentials\|API rate limit exceeded"; then
        echo "Warning: GitHub authentication failed or rate limited, falling back to unauthenticated request"
        response=""
      fi
    fi

    # Fall back to unauthenticated request if no token or authentication failed
    if [ -z "$response" ]; then
      if [ -z "$GITHUB_TOKEN" ]; then
        echo "Using unauthenticated GitHub API request (rate limited to 60 requests/hour)"
      fi
      if type "curl" > /dev/null; then
        response=$(curl -sL "$latest_url")
      elif type "wget" > /dev/null; then
        response=$(wget -q -O - "$latest_url")
      fi
    fi

    # Extract the download URL from the response
    DOWNLOAD_URL=$(echo "$response" | grep "browser_download_url" | grep "${OS}_${ARCH}" | awk '/"browser_download_url":/{gsub( /[,"]/,"", $2); print $2}' | head -n1)

    # Validate that a download URL was found
    if [ -z "$DOWNLOAD_URL" ]; then
      echo "Error: Could not find download URL for ${OS}-${ARCH}"
      echo "Please check that a release exists at: $latest_url"
      exit 1
    fi
  else
    DOWNLOAD_URL="https://github.com/$PROJECT_GH/releases/download/$VERSION/helm-ssm-$OS.tgz"
  fi
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  PLUGIN_TMP_FILE="/tmp/${PROJECT_NAME}.tgz"
  echo "Downloading $DOWNLOAD_URL"
  if type "curl" > /dev/null; then
    curl -L "$DOWNLOAD_URL" -o "$PLUGIN_TMP_FILE"
  elif type "wget" > /dev/null; then
    wget -q -O "$PLUGIN_TMP_FILE" "$DOWNLOAD_URL"
  fi
}

# installFile verifies the SHA256 for the file, then unpacks and
# installs it.
installFile() {
  HELM_TMP="/tmp/$PROJECT_NAME"
  mkdir -p "$HELM_TMP"
  tar xf "$PLUGIN_TMP_FILE" -C "$HELM_TMP"
  echo "$HELM_TMP"
  HELM_TMP_BIN="$HELM_TMP/helm-ssm"
  echo "Preparing to install into ${HELM_PLUGINS}"
  # Use * to also copy the file withe the exe suffix on Windows
  cp "$HELM_TMP_BIN" "$HELM_PLUGINS/helm-ssm"
}

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    echo "Failed to install $PROJECT_NAME"
    echo "For support, go to https://github.com/tutti-ch/helm-ssm."
  fi
  exit $result
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
  set +e
  echo "$PROJECT_NAME installed into $HELM_PLUGINS/$PROJECT_NAME"
  # To avoid to keep track of the Windows suffix,
  # call the plugin assuming it is in the PATH
  PATH=$PATH:$HELM_PLUGINS/$PROJECT_NAME
  helm-ssm -h
  set -e
}

# Execution

#Stop execution on any error
trap "fail_trap" EXIT
set -e
initArch
initOS
verifySupported
getDownloadURL
downloadFile
installFile
testVersion
