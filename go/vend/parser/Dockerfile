FROM saichler/builder:latest AS build
COPY main.go /home/src/github.com/saichler/build/
RUN go mod init
RUN GOPROXY=direct GOPRIVATE=github.com go mod tidy
RUN go build -o vendmachine-parser
FROM saichler/vending-security:latest AS final
COPY --from=build /home/src/github.com/saichler/build/vendmachine-parser /home/run/vendmachine-parser
ENTRYPOINT ["/home/run/vendmachine-parser"]
