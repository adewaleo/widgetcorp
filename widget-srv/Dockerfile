# ---- build stage ----
# Full Go toolchain, only used to compile. Named "builder" so the runtime
# stage can copy the finished binary out of it.
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency manifests first and download modules. Docker caches this
# layer, so deps are only re-downloaded when go.mod/go.sum change — not on
# every source edit.
COPY go.mod go.sum ./
RUN go mod download

# Now the source. (embedded migrations come along via the //go:embed in
# internal/db, so we don't copy .sql files separately.)
COPY . .

# CGO_ENABLED=0 -> a static binary with no libc dependency, so it runs on a
# bare base image. GOOS=linux to be explicit about the target.
RUN CGO_ENABLED=0 GOOS=linux go build -o /widget-srv .

# ---- runtime stage ----
# Tiny base: just enough to run the static binary. Nothing from the toolchain
# above ships in the final image.
FROM alpine:3.20

# Run as a non-root user rather than root.
RUN adduser -D -u 10001 appuser
USER appuser

# Pull only the compiled binary out of the builder stage.
COPY --from=builder /widget-srv /widget-srv

EXPOSE 8080
ENTRYPOINT ["/widget-srv"]
