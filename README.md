# go-allrecipes
🍳🌭 AllRecipes scraper made in one day powered by FastHTTP and Redis 

Created to learn FastHTTP and scrape. Requests cached with Redis.
**Development time so far: 5 hours**

## :cloud: Installation (Linux)

```sh
cd $GOPATH/src
git clone https://github.com/matlsn/go-allrecipes.git
go get
go build
pm2 start go-allrecipes
```

GET: `http://localhost:5557/recipes/<query>`

Ensure [Go is installed properly](https://golang.org/doc/install).

## :dizzy: Where is this API used?

If you are using this API in one of your projects, add it in this list. :sparkles:
