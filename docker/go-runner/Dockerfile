FROM golang:1.22-alpine

WORKDIR /app

RUN adduser -D runner \
    && chown -R runner:runner /app

USER runner

CMD ["sh", "-c", "while true; do sleep 1; done"]