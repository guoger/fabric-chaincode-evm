#
# Copyright IBM Corp. All Rights Reserved
# SPDX-License-Identifier: Apache2.0
#
FROM hyperledger/fabric-buildenv:latest
WORKDIR $GOPATH/src/github.com/hyperledger/
RUN git clone https://github.com/hyperledger/fabric.git
WORKDIR $GOPATH/src/github.com/hyperledger/fabric
RUN git checkout master
RUN EXECUTABLES=go GO_TAGS=pluginsenabled EXPERIMENTAL=false DOCKER_DYNAMIC_LINK=true make peer

FROM hyperledger/fabric-peer:latest
RUN apt-get update -qqy && apt-get dist-upgrade -qqy && apt-get install libltdl-dev -qqy
COPY --from=0 /opt/gopath/src/github.com/hyperledger/fabric/.build/bin/peer /usr/local/bin/peer
RUN ls -l /usr/local/bin
COPY .build/linux/lib/evmscc.so /opt/lib/evmscc.so
COPY config/core.yaml /etc/hyperledger/fabric
