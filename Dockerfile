FROM golang:1.15-alpine as windows-build
ARG GOOS=windows
ARG GOARCH=amd64
ARG OUTPUTFILE=ampt.exe
ARG CGO_ENABLED=1
ARG CC=/usr/bin/x86_64-w64-mingw32-gcc
WORKDIR /src
COPY . .
RUN apk add mingw-w64-gcc && go build -o /out/$GOOS/$GOARCH/$OUTPUTFILE .

FROM scratch as bin
COPY --from=windows-build /out /