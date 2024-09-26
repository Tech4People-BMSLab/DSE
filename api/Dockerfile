########################################
# Alpine
########################################
FROM alpine:latest

RUN apk add ca-certificates
RUN apk add tzdata
RUN apk add bash
RUN apk add git

WORKDIR /app

COPY ./dist/main  ./main
COPY ./data/      ./data/
COPY ./public/    ./public/
COPY ./build.txt  ./build.txt
 
CMD ["./main"]
