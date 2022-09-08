## Introduction
![ggrok-flow](https://github.com/onyas/ggrok/blob/main/docs/flow.jpg?raw=true)

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
or you can deploy with Heroku with `Deploy To Heroku` button below:

[![Deploy To Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

### Client side

Go to the [release page](https://github.com/onyas/ggrok/releases) to download the installation package. After installation, execute the following command:

```
ggrok -proxyServer yourProxyServer.herokuapp.com

ggrok -port 3000
```

Now your local server is exposed to the internet, you can visit it by https://yourProxyServer.herokuapp.com


## Useful docs
- [Architecture of ggrok](./docs/architecture.md)
- [Introduction about ggrok](./docs/introduce.md)
- [How to run this project locally](./docs/run-locally.md)
- [How to debug concurrent issues](./docs/debug.md)
- [Useful Heroku command](./docs/heroku-command.md)
- [How to publish a new release](./docs/new-release.md)
