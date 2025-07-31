#!/bin/bash
set -e

echo "--- Testing auth package ---"
cd /app/module/authentication/auth
go test -v

echo "--- Testing middleware package ---"
cd /app/module/authentication/middleware
go test -v

echo "--- All tests passed ---"
