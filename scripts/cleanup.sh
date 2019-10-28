#!/bin/bash

ctr task kill quake3s
sleep 1
ctr container rm quake3s
ctr snapshot rm quake3s-snapshot
