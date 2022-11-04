#!/usr/bin/env bash

set -xeu

mockery --case underscore --dir ./internal/riot --name Client --output ./internal/riot/mocks
mockery --case underscore --dir ./internal/lol --name Tactics --output ./internal/lol/mocks
