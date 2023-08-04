FROM alpine:latest

WORKDIR /app

COPY app.env /app
COPY accountApp /app

CMD ["/app/accountApp"]