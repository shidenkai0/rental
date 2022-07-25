#!/usr/bin/env bash

doctl registry kubernetes-manifest | kubectl apply -f -
