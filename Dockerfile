FROM golang:1.25 AS build-env
ENV CGO_ENABLED=0
ARG VERSION
ARG GIT_COMMITSHA
ARG GIT_VERSION
ARG GIT_STRIPPED_VERSION

WORKDIR /github.com/meshery/meshery-linkerd
COPY go.mod go.sum ./
RUN go mod download
COPY main.go main.go
COPY internal/ internal/
COPY linkerd/ linkerd/
COPY build/ build/
RUN VERSION=$(curl -L -s \
    https://api.github.com/repos/meshery/meshery-linkerd/releases/latest | \
	grep tag_name | sed "s/ *\"tag_name\": *\"\\(.*\\)\",*/\\1/" | \
	grep -v "rc\.[0-9]$"| head -n 1 )
RUN go mod tidy; CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build \
	-ldflags="-w -s -X main.version=$VERSION -X main.gitsha=$GIT_COMMITSHA" \
	-a -o meshery-linkerd main.go



# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/nodejs:16
ENV DISTRO="debian"
ENV SERVICE_ADDR="meshery-linkerd"
ENV MESHERY_SERVER="http://meshery:9081"
WORKDIR /
COPY templates/ ./templates
COPY --from=build-env /github.com/meshery/meshery-linkerd/meshery-linkerd .
ENTRYPOINT ["./meshery-linkerd"]
