FROM golang:latest

COPY main.go /app/main.go
RUN go get github.com/labstack/echo && go get go.mongodb.org/mongo-driver/mongo && go get github.com/dgrijalva/jwt-go

CMD ["go", "run", "/app/main.go"]
