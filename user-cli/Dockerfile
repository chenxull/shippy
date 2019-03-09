FROM alpine:latest
RUN mkdir /app
WORKDIR /app
ADD user-cli /app/user-cli
ENTRYPOINT ["./user-cli"]
CMD ["user-cli"]