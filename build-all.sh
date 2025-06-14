#!/bin/bash

gox -os="darwin" -os="linux" -os="windows" -arch="amd64" -arch="arm64" -osarch="linux/386" -osarch="windows/386"
