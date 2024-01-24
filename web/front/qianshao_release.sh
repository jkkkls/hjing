#!/bin/bash

rm -rf dist dist.zip

yarn build
zip -qr dist.zip ./dist

scp ./dist.zip root@xxx:~/update
ssh root@xxx "cd ~/update;pwd;./update_xxoo.sh"
