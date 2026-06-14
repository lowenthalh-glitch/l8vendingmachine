FROM saichler/builder:latest AS build

COPY main.go /home/src/github.com/saichler/build/main.go
RUN go mod init
RUN GOPROXY=direct GOPRIVATE=github.com go mod tidy
RUN go build -o logagent

FROM saichler/vending-security:latest AS final
COPY --from=build /home/src/github.com/saichler/build/logagent /home/run/logagent

ENTRYPOINT ["/home/run/logagent"]
