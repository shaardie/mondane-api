FROM golang AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY main.go ./
COPY api ./api
RUN go build -o mondane-api

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=build /app/mondane-api /mondane-api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/mondane-api"]
