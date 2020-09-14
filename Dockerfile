FROM alpine:3
WORKDIR /
COPY dmicrog /
ENTRYPOINT ["/dmicrog"]
CMD ["$@"]