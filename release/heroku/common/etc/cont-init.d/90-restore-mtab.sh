#!/bin/bash

# fix the missing `/etc/mtab` file error
if [[ ! -f /etc/mtab ]]; then
    ln -sf /proc/mounts /etc/mtab
fi