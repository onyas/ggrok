## Usage

```
go run main.go

go run main.go -client
```

request http://localhost:8080/ through postman



[Data Race Detector](https://go.dev/doc/articles/race_detector)

```
go run -race main.go

go run -race main.go -client
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
```