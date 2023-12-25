#!/bin/bash

# Step 1: Build your Go binary
make build
make run

pm2 save
pm2 startup

