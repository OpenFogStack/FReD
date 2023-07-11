# Base image: https://hub.docker.com/_/golang/
FROM golang:1.20

ENV DOCKER_CHANNEL stable
ENV DOCKER_VERSION 24.0.2
ENV DOCKER_BUILDX_VERSION 0.11.1
ENV GOLANGCI_LINT_VERSION v1.53.3

# Install golint
ENV GOPATH /go
ENV PATH ${GOPATH}/bin:$PATH
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

# Install clang from LLVM repository
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install Docker (to build images)

RUN set -eux; \
    \
    apkArch="$(arch)"; \
    case "$apkArch" in \
    # amd64
    x86_64) dockerArch='x86_64' ;; \
    # arm32v6
    armhf) dockerArch='armel' ;; \
    # arm32v7
    armv7) dockerArch='armhf' ;; \
    # arm64v8
    aarch64) dockerArch='aarch64' ;; \
    *) echo >&2 "error: unsupported architecture ($apkArch)"; exit 1 ;;\
    esac; \
    \
    if ! wget -O docker.tgz "https://download.docker.com/linux/static/${DOCKER_CHANNEL}/${dockerArch}/docker-${DOCKER_VERSION}.tgz"; then \
    echo >&2 "error: failed to download 'docker-${DOCKER_VERSION}' from '${DOCKER_CHANNEL}' for '${dockerArch}'"; \
    exit 1; \
    fi; \
    \
    tar xzvf docker.tgz --strip=1 -C /usr/local/bin docker/docker; \
    rm docker.tgz; \
    \
    docker --version

# Install docker buildx

RUN set -eux; \
    \
    apkArch="$(arch)"; \
    case "$apkArch" in \
    # amd64
    x86_64) dockerArch='amd64' ;; \
    # arm32v6
    armhf) dockerArch='arm-v6' ;; \
    # arm32v7
    armv7) dockerArch='arm-v7' ;; \
    # arm64v8
    aarch64) dockerArch='arm64' ;; \
    *) echo >&2 "error: unsupported architecture ($apkArch)"; exit 1 ;;\
    esac; \
    \
    if ! wget -O docker-buildx "https://github.com/docker/buildx/releases/download/v${DOCKER_BUILDX_VERSION}/buildx-v${DOCKER_BUILDX_VERSION}.linux-${dockerArch}"; then \
    echo >&2 "error: failed to download 'buildx-${DOCKER_BUILDX_VERSION}' for '${dockerArch}'"; \
    exit 1; \
    fi; \
    \
    mkdir -p /usr/local/lib/docker/cli-plugins; \
    mv docker-buildx /usr/local/lib/docker/cli-plugins; \
    chmod +x /usr/local/lib/docker/cli-plugins/docker-buildx; \
    docker buildx version
