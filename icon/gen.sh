#/bin/sh

if [ -z "$GOPATH" ]; then
  echo GOPATH environment variable not set
  exit
fi

if [ ! -e "$GOPATH/bin/2goarray" ]; then
  echo "Installing 2goarray..."
  go get github.com/cratonica/2goarray
  if [ $? -ne 0 ]; then
    echo Failure executing go get github.com/cratonica/2goarray
    exit
  fi
fi

generate_icon() {
  local IMGPATH="$1" # path to image
  local OUTPUT="$2"  # output .go file name
  local GONAME="$3"  # go prop name
  local TARGET="$4"  # unix or win

  BUILD="//+build linux darwin"
  if [ $TARGET == "win" ]; then
    BUILD="//+build windows"
  fi

  echo Generating $OUTPUT
  echo "$BUILD" > $OUTPUT
  echo "//goland:noinspection ALL" >> $OUTPUT
  echo >> $OUTPUT
  cat "$IMGPATH" | $GOPATH/bin/2goarray $GONAME icon >> $OUTPUT
  if [ $? -ne 0 ]; then
    echo Failure generating $OUTPUT
  fi
}

generate_icon "./assets/base.png" icon_base.go baseIcon unix
# generate_icon "./assets/base.ico" icon_base_win.go baseIcon win

echo Finished