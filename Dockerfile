FROM golang:1.16-alpine
ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ARG GITHUB_CLIENT_ID
ARG GITHUB_CLIENT_SECRET
ARG SESSION_KEY
ENV GITHUB_CLIENT_ID $GITHUB_CLIENT_ID
ENV GITHUB_CLIENT_SECRET $GITHUB_CLIENT_SECRET
ENV SESSION_KEY $SESSION_KEY
RUN mkdir -p /build
WORKDIR /build
COPY go.* /build/
RUN go mod download
COPY . .
RUN go build -o go-github-tenable
CMD ["./go-github-tenable"]
EXPOSE 8080