#!/bin/bash

kill -SIGUSR2 $(cat httpdns.pid)
