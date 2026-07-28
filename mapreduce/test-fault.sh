#!/bin/bash

echo "building binary..."
go build -o mr-app .
if [ $? -ne 0 ]; then
    echo "build failed"
    exit 1
fi

rm -f temp/* out/*

echo "starting master"
./mr-app -role=master &
MASTER_PID=$!

sleep 1

echo "starting workers"
./mr-app -role=worker &
WORKER1=$!
./mr-app -role=worker &
WORKER2=$!
./mr-app -role=worker &
WORKER3=$!

sleep 2

kill -9 $WORKER2

sleep 2

echo "replacement worker"
./mr-app -role=worker &
WORKER4=$!

wait $MASTER_PID
wait $WORKER1
wait $WORKER2
wait $WORKER3
wait $WORKER4
echo "done"
ls -lh out/
