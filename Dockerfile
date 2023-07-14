FROM golang as builder

WORKDIR /app
COPY . .
RUN go build -ldflags "-w"

FROM ubuntu:22.04

COPY --from=builder /app/syslog-perf /usr/bin/syslog-perf

ARG TZ=UTC
RUN set -ex;                                         \
  ln -snf /usr/share/zoneinfo/$TZ /etc/localtime;    \
  apt-get update;                                    \
  apt-get install --no-install-recommends -y         \
    busybox                                          \
    pid1                                             \
    tzdata                                           \
  ;                                                  \
  rm -rf /var/lib/apt/lists/*;                       \
  busybox --install

ENTRYPOINT [ "pid1", "--" ]
CMD [ "tail", "-f", "/dev/null" ]
