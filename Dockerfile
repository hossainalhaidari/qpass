FROM golang:alpine AS build
WORKDIR /src
COPY . /src/
RUN go build -o /bin/qpass

FROM scratch
WORKDIR /app
COPY --from=build /bin/qpass ./qpass
COPY --from=build /src/index.html ./index.html
EXPOSE 3000
ENTRYPOINT ["/app/qpass"]
