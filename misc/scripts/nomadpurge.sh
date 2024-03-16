#!/bin/zsh
nomad job stop -purge pytechco-employee
nomad job stop -purge pytechco-web
nomad job stop -purge pytechco-redis
nomad job stop -purge pytechco-setup