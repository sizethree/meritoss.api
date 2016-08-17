#!/usr/bin/env bash

set +x

function distro {
  local UNAME=`uname`
  local DISTRO=""

  if [ $UNAME == "Linux" ]; then
    DISTRO="linux-amd64"
  fi

  if [ $UNAME == "Darwin" ]; then
    DISTRO="darwin-amd64"
  fi

  echo $DISTRO
}

function version {
  cat <<EOM
Golang Environment:
  go version: $(go version)
  GOPATH:     $GOPATH
  GOROOT:     $GOROOT
  PATH:       $PATH
  HOME:       $HOME
EOM
}

function upgrade {
  local VERSION=1.6.3
  local DISTRO=$(distro)

  local MERITOSS_ROOT=$(echo $HOME)/.meritoss

  if [ ! -z $APP_RUNTIME_DIR ]; then
    printf "Using app runtime dir: $APP_RUNTIME_DIR\n"
    MERITOSS_ROOT=$APP_RUNTIME_DIR
  fi

  local DOWNLOAD_DIR=./build/go-$VERSION
  local INSTALL_DIR=$MERITOSS_ROOT/go-$VERSION

  if [ -z $DISTRO ]; then
    printf "Unable to determine distro file from \"`uname`\", exiting.\n"
    exit 1
  fi

  # print our current version info
  version

  # generate the string that represents the artifact file
  local ARTIFACT=go$VERSION.$DISTRO.tar.gz

  # create a directory to download into
  mkdir -p $DOWNLOAD_DIR
  
  printf "\nDownloading go v$VERSION for $DISTRO...\n"
  # download the file
  curl -s -O "https://storage.googleapis.com/golang/$ARTIFACT"

  # extract into download dir
  tar -zxf $ARTIFACT -C $DOWNLOAD_DIR

  # remove the tarball
  rm $ARTIFACT

  rm -rf $INSTALL_DIR
  mkdir -p $INSTALL_DIR

  # prepare our root and bin path env vars
  export GOROOT=$INSTALL_DIR
  export GOPATH=$GOROOT/bin

  rm -rf $GOROOT
  mv $DOWNLOAD_DIR/go $GOROOT

  export PATH=$GOPATH:$PATH

  printf "Installed go v$VERSION into $INSTALL_DIR, new environment information:\n"
  rm -rf $DOWNLOAD_DIR
}

function install {
  local DEST=$GOROOT/src/github.com/sizethree/meritoss.api
  local EXE=meritoss.api
  local CWD=$(pwd)

  printf "Installing project into $DEST\n"
  version

  rm -rf $DEST
  mkdir -p $DEST
  cp -r * $DEST

  printf "Installing dependencies...\n"
  go get -u github.com/kardianos/govendor

  cd $DEST
  govendor sync +e,^local
  go build -o $EXE
  chmod +x $EXE

  printf "Compilation complete, installing to $CWD/$EXE\n"
  cd $CWD
  mv $DEST/$EXE ./
}

if [ -z $1 ]; then
  upgrade
  install
else
  echo $1
fi