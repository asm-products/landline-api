# Landline API

<a href="https://assembly.com/landline/bounties?utm_campaign=assemblage&utm_source=landline&utm_medium=repo_badge"><img src="https://asm-badger.herokuapp.com/landline/badges/tasks.svg" height="24px" alt="Open Tasks" /></a>
[![Build Status](https://travis-ci.org/asm-products/landline-api.png?branch=master)](https://travis-ci.org/asm-products/landline-api)

## Drop in chat for your app

This is a product being built by the Assembly community. You can help push this idea forward by visiting [https://assembly.com/landline](https://assembly.com/landline).

## Development

The easiest way to run the API locally is with [Compose](https://docs.docker.com/compose/). First make sure Docker and Compose are installed. Then run:

    ./dc-setup
    ./db/dc-migrate
    ./db/dc-test-data
    docker-compose run web go run db/test-data.go
    docker-compose up

If you want to run it outside Docker, make sure postgres and go version 1.4 are installed, then run    

    go get bitbucket.org/liamstask/goose/cmd/goose
    go get github.com/codegangsta/gin
    godep restore
    ./db/setup
    ./db/migrate
    forego run go run db/test-data.go
    gin

### Authentication flow

    # Start up the test id provider
    $ go run example/identity_provider.go 41fe7589256fd058b3f56bc71a56ebad3b1d6b86e027a73a02db0e3a0524f9d4

    # Grab the sso endpoint with nonce from landline
    $ curl -i 'localhost:3000/sessions/new?team=test-dev'
    HTTP/1.1 302 Found
    Location: http://localhost:8989/sso?payload=...&sig=....

    # Hit the identity provider with payload landline redirect, this needs to return a user payload
    $ curl -i http://localhost:8989/sso?payload=...&sig=....
    HTTP/1.1 302 Found
    Location: http://localhost:3000/sessions/sso?payload=...&sig=...

    # hit landline with the user payload to receive a session token for the api
    $ curl -i http://localhost:3000/sessions/sso?payload=...&sig=...
    {"token":"..."}

    # now you can make requests with your jwt token
    $ curl -H "Authorization: Bearer $TOKEN" localhost:3000/rooms
    {"rooms":[]}

### Easily test the api using JavaScript

If you're running the example identity provider at port 8989, you can go to localhost:8989/debug. You'll find an empty page, but when you open up the javascript console in your browser of choice, you'll see that there are a bunch of handy Javascript objects and functions to help you debug your changes to the api.

#### Promises ####

The javascript on this page uses Promises. If you're not familiar with them, check out [this html5rocks article on them](http://www.html5rocks.com/en/tutorials/es6/promises/), they're awesome.

#### The  Session object


The session object helps you make authenticated calls to the API. You don't construct them using the `new` keyword, but by calling `Session.create()`. This function returns a promise, which will resolve to a session as soon as we've obtained a session token from the API server.
```js
Session.create(function(sess){
    // You can use 'sess' as the session in here.
});
```

#### Session.immediate()

When you're writing scripts, it's handy to know when your session token has been obtained, and your session is ready to use. But when you're just testing stuff in the javascript console, the request will probably be completed multiple seconds before you're done typing your next command. Session.immediate allows you to treat session creation as if it's asynchronous.
```js
sess = Session.immediate();
// You can 'immediately' start using the session object.
```

#### session.makeCall(method, path, data)

On a session object, you can call session.makeCall to make an api call. `method` is the HTTP request method, path is the path you'd like to request, including query parameters. If `data` is given, it'll be json serialized, and sent as the request body. This method returns a promise, which will resolve to the JSON parsed response text.

```js
Session.create(function(sess){
    sess.makeCall("GET", "/rooms").then(roomsResponse){
        // use the response in here.
    }
});

```

#### boundLog

the `boundLog` function does exactly the same thing as `console.log`, but it is already bound to the console object, so you can easily pass it as a callback.

```js
//within the javascript console:
sess = Session.immediate();
sess.makeCall("GET", "/rooms").then(boundLog);
```

This will log the javascript object returned by the /rooms endpoint to the console.

### How Assembly Works

Assembly products are like open-source and made with contributions from the community. Assembly handles the boring stuff like hosting, support, financing, legal, etc. Once the product launches we collect the revenue and split the profits amongst the contributors.

Visit [https://assembly.com](https://assembly.com) to learn more.
