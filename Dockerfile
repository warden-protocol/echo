FROM gcr.io/distroless/static-debian12
# FROM ubuntu:22.04

COPY ./echo /

EXPOSE 10010

CMD ["/echo"]
