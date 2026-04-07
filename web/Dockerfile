# Stage 1: Build WASM binary + copy wasm_exec.js
FROM golang:1.26-alpine AS go-builder
WORKDIR /src
COPY . .
RUN cd wasm && GOOS=js GOARCH=wasm go build -ldflags="-s -w" \
    -o ../web/app/public/passgen.wasm .
# Try both known wasm_exec.js locations (path changed across Go versions)
RUN cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" web/app/public/wasm_exec.js 2>/dev/null || \
    cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/app/public/wasm_exec.js

# Stage 2: Build Svelte app
FROM node:24-alpine AS node-builder
WORKDIR /app
COPY web/app/package*.json ./
RUN npm ci
COPY --from=go-builder /src/web/app/public/ ./public/
COPY web/app/ .
RUN npm run build

# Stage 3: Build Go web server
FROM golang:1.26-alpine AS server-builder
WORKDIR /src
COPY . .
COPY --from=node-builder /app/dist web/app/dist/
# Build from within the web module so the replace directive (../) resolves correctly
RUN cd web && CGO_ENABLED=0 go build -ldflags="-s -w" \
    -o ../passgen-web ./cmd/passgen-web/

# Final: minimal scratch image
FROM scratch
COPY --from=server-builder /src/passgen-web /passgen-web
EXPOSE 8080
ENTRYPOINT ["/passgen-web"]
