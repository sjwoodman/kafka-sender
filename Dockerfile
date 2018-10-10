FROM centos:7
MAINTAINER Simon Woodmna <swoodman@redhat.com>
ARG BINARY=./kafkasender

COPY ${BINARY} /opt/kafkasender
ENTRYPOINT ["/opt/kafkasender"]