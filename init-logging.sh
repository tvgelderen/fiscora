#!/bin/bash

LOG_DIR_FRONTEND="/var/log/fiscora-frontend"
LOG_DIR_BACKEND="/var/log/fiscora-backend"

mkdir -p "$LOG_DIR_FRONTEND"
mkdir -p "$LOG_DIR_BACKEND"

chmod -R 777 "$LOG_DIR_FRONTEND"
chmod -R 777 "$LOG_DIR_BACKEND"
