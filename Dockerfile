FROM oven/bun as js_build

COPY ./ /app/
WORKDIR /app

RUN bun install
RUN bun run build

FROM golang:1.20-alpine AS base
FROM base AS build

ENV PORT=8080

COPY ./ /app/
WORKDIR /app

RUN go mod download

COPY --from=js_build /app/public/ /app/public/
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

FROM base

ENV PORT=8080

COPY --from=build /app/main /app/main

EXPOSE 8080

CMD ["/app/main"]
