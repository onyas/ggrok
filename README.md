## Usage

### Server side
1. Fork and clone this repo
2. Deploy to heroku
    ```
    //create app in heroku
    heroku create
    //deploy your app to heroku
    git push heroku main
    ```

### Client side
```
go run main.go -client -proxyServer yourProxyServer.herokuapp.com

go run main.go -client -port 3000
```

Now your local server is exposed to the internet, you could visit it by http://yourProxyServer.herokuapp.com


## Useful docs
- [Architecture of ggrok](./docs/architecture.md)
- [How to run this project locally](./docs/run-locally.md)
- [How to debug concurrent issues](./docs/debug.md)
- [Useful Heroku command](./docs/heroku-command.md)