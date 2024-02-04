#!/usr/bin/env bash
CRed=$(tput setaf 1)
CBlue=$(tput setaf 4)
CCyan=$(tput setaf 6)
CGreen=$(tput setaf 2)
CYellow=$(tput setaf 3)
BOLD=$(tput smso)
UNBOLD=$(tput rmso)
ENDMARKER=$(tput sgr0)

repo=smash
name=smash

echo "----------------------"
echo "${CGreen}${name} Installer v1.0.0${ENDMARKER}"
echo "----------------------"


function fatal() {
    echo ""
    echo "${BOLD}${CRed}FATAL:${UNBOLD}${ENDMARKER} $1${ENDMARKER}" >&2
    echo ""
    exit 1
}

function extract_zip() {
    local filename=$1
    if [[ -x "$(command -v unzip)" ]]; then
        echo "${CCyan}RUN${ENDMARKER} unzip -xvf ${CYellow}'${filename}'${ENDMARKER}"
        unzip -o "${filename}" -d "${filename%.*}/"
        if [ $? -eq 0 ]; then
            echo "${CCyan}RAN${ENDMARKER} Extracted ${CGreen}smash ${latest_release}${ENDMARKER} to ${CYellow}'${filename}/'${ENDMARKER}"
            rm -f "${filename}"
            echo "${CCyan}DEL${ENDMARKER} rm ${CYellow}'${filename}'${ENDMARKER}"
        else
            fatal "Failed to extract ${filename}"
        fi
    else
        echo "${BOLD}${CRed}ERROR:${UNBOLD}${ENDMARKER} unzip not found, please extract manually!${ENDMARKER}" >&2
    fi
}
function extract_tar() {
    local filename=$1
    if [[ -x "$(command -v tar)" ]]; then
        echo "${CCyan}RUN${ENDMARKER} tar -xvf ${CYellow}'${filename}'${ENDMARKER}"
        tar -xvf "${filename}" --one-top-level
        if [ $? -eq 0 ]; then
            echo "${CCyan}RAN${ENDMARKER} Extracted ${CGreen}smash ${latest_release}${ENDMARKER} to ${CYellow}'${filename}/'${ENDMARKER}"
            rm -f "${filename}"
            echo "${CCyan}DEL${ENDMARKER} rm ${CYellow}'${filename}'${ENDMARKER}"
        else
            fatal "Failed to extract ${filename}"
        fi

    else
        echo "${BOLD}${CRed}ERROR:${UNBOLD}${ENDMARKER} tar not found, please extract manually!${ENDMARKER}" >&2
    fi
}

if [[ "$OSTYPE" == "linux-gnu"* ]] || [[ "$OSTYPE" == "linux-musl" ]] ; then
  os="linux"
  ext="tar.gz"
elif [[ "$OSTYPE" == "darwin"* ]]; then
  os="macos"
  ext="zip"
elif [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
  os="windows"
  ext="zip"
elif [[ "$OSTYPE" == "freebsd"* ]]; then
  os="freebsd"
  ext="tar.gz"
else
  fatal "Unsupported OS Type $OSTYPE"
fi

# Detect the architecture
case $(uname -m) in
    x86_64)
        arch="amd64"
        ;;
    arm64)
        arch="arm64"
        ;;
    *)
        fatal "Unsupported architecture, only amd64 and arm64 are supported"
        ;;
esac

echo "OS:     ${OSTYPE}"
echo "ARCH:   ${arch}"

latest_release=$(curl --silent "https://api.github.com/repos/thushan/$repo/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$latest_release" ]; then
    fatal "Failed to get the latest release"
fi
echo "LATEST: ${latest_release}"
echo "----------------------"

# Construct the file name
file="${name}_${latest_release}_${os}_${arch}.${ext}"

# Construct the download URL
url="https://github.com/thushan/$repo/releases/download/${latest_release}/${file}"

echo "${CCyan}GET${ENDMARKER} via ${CBlue}'${url}'${ENDMARKER}"
# Download the release
curl --silent -L "$url" -o "$file"

# Check if the download was successful
if [ $? -eq 0 ]; then
    echo "${CCyan}GOT${ENDMARKER} ${CGreen}smash ${latest_release}${ENDMARKER} Downloaded to ${CYellow}'${file}'${ENDMARKER}"
    if [[ "$ext" == "zip" ]]; then
        extract_zip "$file"
    else
        extract_tar "$file"
    fi
else
    fatal "Failed to download ${url}"
fi

echo "${CCyan}YAY${ENDMARKER} Installed ${CGreen}smash ${latest_release}${ENDMARKER} to ${CYellow}'${file}'${ENDMARKER}"
