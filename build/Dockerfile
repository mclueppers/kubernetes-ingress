# Copyright 2019 HAProxy Technologies LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.12.9-alpine AS builder

RUN apk --no-cache add git openssh

COPY /go.mod /src/go.mod
COPY /go.sum /src/go.sum
RUN cd /src && go mod download

COPY / /src

RUN mkdir -p /var/run/vars && \
    cd /src && \
    git config --get remote.origin.url > /var/run/vars/GIT_REPO && \
    git rev-parse --short HEAD > /var/run/vars/GIT_HEAD_COMMIT && \
    git describe --abbrev=0 --tags > /var/run/vars/GIT_LAST_TAG && \
    git rev-parse --short $(cat /var/run/vars/GIT_LAST_TAG) > /var/run/vars/GIT_TAG_COMMIT && \
    git diff $(cat /var/run/vars/GIT_HEAD_COMMIT) $(cat /var/run/vars/GIT_TAG_COMMIT) --quiet > /var/run/vars/GIT_MODIFIED1 || echo '.dev' > /var/run/vars/GIT_MODIFIED1 && \
    git diff --quiet > /var/run/vars/GIT_MODIFIED2 || echo '.dirty' > /var/run/vars/GIT_MODIFIED2 && \
    cat /var/run/vars/GIT_MODIFIED1 /var/run/vars/GIT_MODIFIED2 | tr -d '\n' > /var/run/vars/GIT_MODIFIED && \
    date '+%Y-%m-%dT%H:%M:%S' > /var/run/vars/BUILD_DATE && \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
        -gcflags "-N -l" \
        -ldflags "-X main.GitRepo=$(cat /var/run/vars/GIT_REPO) -X main.GitTag=$(cat /var/run/vars/GIT_LAST_TAG) -X main.GitCommit=$(cat /var/run/vars/GIT_HEAD_COMMIT) -X main.GitDirty=$(cat /var/run/vars/GIT_MODIFIED) -X main.BuildTime=$(cat /var/run/vars/BUILD_DATE)" \
        -o fs/haproxy-ingress-controller .

FROM mclueppers/haproxy-alpine:v2.0.4

RUN apk --no-cache add socat openssl util-linux openrc htop dumb-init

COPY /fs /
COPY --from=builder /src/fs/haproxy-ingress-controller .

RUN chmod +x /etc/init.d/haproxy

ENTRYPOINT ["/usr/bin/dumb-init", "--", "/start.sh"]
