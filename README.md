## Usage
go run main.go -client -proxyServer https://yourProxyServer.com
go run main.go -client -port 3000


## Debug stage

[Data Race Detector](https://go.dev/doc/articles/race_detector)
```
go run -race main.go

go run -race main.go -client
```
request http://localhost:8080/ through postman

## Development with Heroku locally
```
//start server locally
heroku local -f Procfile.local

//start client locally
go run main.go -client
```

## Deploy to Heroku

```
heroku login
//create app in heroku
heroku create
//deploy your app to heroku
git push heroku main
//open you app in browser
heroku open

//check logs
heroku logs --tail
```