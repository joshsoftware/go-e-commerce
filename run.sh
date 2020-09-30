#!/bin/bash
echo "script started"
echo "Migration started"
./bin/go-e-commerce migrate
echo "Migration Completed"
./bin/go-e-commerce start
echo "Server started"