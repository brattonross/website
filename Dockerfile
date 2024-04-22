FROM oven/bun:latest as client

WORKDIR /app
COPY . .

RUN bun install --frozen-lockfile
RUN make build-client

FROM golang:1.22 as server

WORKDIR /app
COPY . .

RUN make build-server

FROM alpine:3.19

COPY --from=client /app/public/main.js /app/public/main.js
COPY --from=client /app/public/styles.css /app/public/styles.css

COPY --from=server /app/bin/server /app/bin/server
COPY --from=server /app/content/ /app/content/
COPY --from=server /app/templates/ /app/templates/

ENV HOST=0.0.0.0
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["./bin/server"]
