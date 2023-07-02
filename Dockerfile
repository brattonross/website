FROM node:20-slim AS node_build

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

RUN corepack enable

COPY ./ /app/
WORKDIR /app

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

FROM golang:1.20-alpine AS base
FROM base AS build

ENV PORT=8080

COPY ./ /app/
WORKDIR /app

RUN go mod download

COPY --from=node_build /app/public/ /app/public/
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

FROM base

ENV PORT=8080

COPY --from=build /app/main /app/main

EXPOSE 8080

CMD ["/app/main"]
