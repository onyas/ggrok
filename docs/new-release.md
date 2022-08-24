brew install goreleaser

goreleaser init

configure you github token in ~/.config/goreleaser/github_token


git tag -a v0.1.0 -m "release v0.1.0"

git push origin v0.1.0

rm -rf dist/*

goreleaser --rm-dist