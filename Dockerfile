FROM oven/bun as build

COPY ./ /app/
WORKDIR /app

RUN bun install
RUN bun run build

FROM oven/bun as base

WORKDIR /app
COPY --from=build /app/dist/ .

ENV PORT=8080
EXPOSE 8080

CMD ["bun", "run", "./server/entry.mjs"]
