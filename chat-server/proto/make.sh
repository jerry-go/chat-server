#!/bin/sh

pushd . > /dev/null
cd source
PROTO_FILES=""
for file in *.proto
do
	echo $file
	if test -f $file
	then
		PROTO_FILES="$PROTO_FILES $file"
	fi
done
protoc -I . $PROTO_FILES --go_out=plugins=grpc:../
popd > /dev/null
ls *.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
