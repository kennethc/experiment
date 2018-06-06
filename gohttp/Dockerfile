#FROM golang

#WORKDIR /go/src/github.com/kennethc/experiment/gohttp
#COPY . .
#RUN go build -ldflags "-linkmode external -extldflags -static" -a main.go

FROM scratch
#COPY --from=0 /go/src/github.com/kennethc/experiment/gohttp/main /main
COPY main /main
EXPOSE 8080
CMD ["/main"]
