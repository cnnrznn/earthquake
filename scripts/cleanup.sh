#!/bin/bash

ctr task kill quake3s
sleep 1
ctr container rm quake3s
ctr image rm c1
ctr image rm checkpoint
