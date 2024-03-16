#!/bin/zsh
export NOMAD_ADDR=http://localhost:4646
sleep 2
open http://localhost:4646
sleep 2
nomad job run pytechco-redis.nomad.hcl
sleep 2
nomad job run pytechco-web.nomad.hcl
sleep 2
nomad job run pytechco-setup.nomad.hcl
sleep 2
nomad job dispatch -meta budget="200000" pytechco-setup
sleep 2
nomad job run pytechco-employee.nomad.hcl
sleep 2
nomad node status -verbose \
    $(nomad job allocs pytechco-web | grep -i running | awk '{print $2}') | \
    grep -i ip-address | awk -F "=" '{print $2}' | xargs | \
    awk '{print "http://"$1":5000"}'


