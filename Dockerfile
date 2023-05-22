### Description: Dockerfile for opensource-insight-exporter
FROM alpine:3.16

COPY opensource-insight-exporter /

# Starting
ENTRYPOINT [ "/opensource-insight-exporter" ]