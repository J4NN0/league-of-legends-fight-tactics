#!/usr/bin/env bash

set -xeu

mockery --case underscore --dir ./pkg/riot --name Client --output ./pkg/riot/mocks
mockery --case underscore --dir ./pkg/lol --name Tactics --output ./pkg/lol/mocks
