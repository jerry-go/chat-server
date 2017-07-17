#!/bin/sh
if [ $# -eq 0 ]; then
    echo "miss GOOS(linux/darwin)"
    exit 1
fi

GOOS=$1
if [[ "$GOOS" != "linux" && "$GOOS" != "darwin" ]]; then
    echo "invalid GOOS $GOOS"
    exit 1
fi

TARGETS=`ls ./apps`
if [ $# -gt 1 ]; then
	args=($*)   
	TARGETS=(${args[@]:1:$#})   
fi

BIN=`pwd`/bin/$GOOS
if [ ! -d $BIN ]; then
    mkdir -p $BIN
fi

echo "GOOS: ${GOOS}"
for target in ${TARGETS[@]}
do
    printf "%s" $target
    pushd . > /dev/null
    cd "./apps/$target"
    OUT="${target}.$GOOS"
    GOOS=$GOOS GOARCH=amd64 go build -o "${BIN}/$OUT"
    popd > /dev/null
    printf " --> %s\n" $OUT
done
