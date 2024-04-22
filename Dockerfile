FROM oven/bun:latest as client

WORKDIR /app
COPY . .

RUN bun install --frozen-lockfile
RUN bun build main.js --outdir public
RUN bun tailwindcss -i styles.css -o public/styles.css

FROM golang:1.22 as server

WORKDIR /app
COPY . .

RUN go build -o ./bin/server ./main.go

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
