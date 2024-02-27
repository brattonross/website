FROM oven/bun:latest as client

WORKDIR /app
COPY . .

RUN bun install --frozen-lockfile
RUN bun run build

FROM golang:1.22 as server

WORKDIR /app
COPY . .

RUN go build -o server

COPY --from=client /app/public/main.js /app/public/main.js
COPY --from=client /app/public/styles.css /app/public/styles.css

ENV HOST=0.0.0.0
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["./server"]
