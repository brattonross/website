FROM golang:1.22 as build

WORKDIR /app

RUN apt-get update && apt-get install -y unzip
RUN curl -fsSL https://bun.sh/install | BUN_INSTALL=/usr bash

COPY . .

RUN bun install --frozen-lockfile
RUN make build

FROM alpine:3

WORKDIR /app

COPY --from=build /app/bin ./bin
COPY --from=build /app/content ./content
COPY --from=build /app/public ./public
COPY --from=build /app/templates ./templates

RUN ["ls", "/app"]
RUN ["ls", "/app/bin"]
RUN ["ls", "."]
RUN ["ls", "./bin"]

ENV HOST=0.0.0.0
ENV PORT=8080
EXPOSE 8080
CMD ["/app/bin/server"]
