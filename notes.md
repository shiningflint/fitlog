## Go HTMX Examples
The simplest type with go htmx and bootstrap
https://github.com/danielmeint/go-htmx-tasklist

This one uses a small framework: leapkit
I think this one is more organized. I will keep an eye for the foldering organization
https://github.com/paganotoni/todox

leapkit template
https://socket.dev/go/package/github.com/leapkit/template?section=overview

This one just uses node, not recommended
https://github.com/Jason-CKY/go-htmx-example/tree/main

## How I write HTTP services in Go after 13 years
https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/

## Go Web App Structure
https://www.reddit.com/r/golang/comments/ttwole/anyone_here_using_go_for_more_traditional_web/

 I've been writing Go professionally for the last 5 years, and we've almost exclusively been using it for backend web development, both as monolithic applications and service-oriented architectures. The standard library covers 99% of our needs. It's really not super complicated, imo. We don't really use gin or any of the other frameworks, as we don't really think they're necessary. Anyone with Go experience already knows all of the things we're already doing.

I don't have any code samples to share with you since everything is proprietary, but the gist of how we set things up is essentially something like this:
- Cobra for CLI commands (serve, migrate, etc)
- gorilla/mux for routing

 Directory structure looks something like:
```go
/cmd
    root.go
    serve.go
/internal
    /sample
        controller.go
        requests.go
        responses.go
main.go
```

The serve.go file bootstraps the app and looks something along the lines of:
```go
router := mux.NewRouter()

controller := sample.Controller{
   Logger: myLogger, // whatever logging package you like
}

controller.Register(router) // register routes

log.Println("listening on port 8080")
log.Fatalln(http.ListenAndServe(":8080", router))
```

and then the controller utilizes a struct so we can inject dependencies. Looks something like this:
```go
package sample

... // imports go here, omitted for brevity

type Controller struct {
    Logger log.Logger
}

func (c Controller) Register(router *mux.Router) {
    router.HandleFunc("/v1/hello", c.Hello()).Methods(http.MethodGet)
}

func (c Controller) Hello() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        c.Logger.Info("hello request received")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{ "greeting": "Hello World" }`))
    }
}
```

The requests.go would contain any structs that represent requests for any of the controllers, particularly POST/PUT bodies in JSON. The responses.go controller represents any responses. For the above example, we might have a response as such:

```go
type Hello struct {
    Greeting string `json:"greeting"`
}
```

and then in the controller, do something like:

```go
resp := response.Hello{ Greeting: "Hello World" }
j, _ := json.Marshal(resp) // err check skipped for brevity
w.Write(j)
```

We tend to prefer the above patterns for a number of reasons that probably aren't worth going into right now, but this tends to work fairly well for us and we have had very few problems. We feel this is fairly idiomatic Go, as well. A lot of folks coming from Java or Rails or Django backgrounds tend to structure their code more in the MVC-style or abstract things a bit more than the above example suggests. It's sometimes very obvious when a developer is coming from another language, particularly an OOP one, and is trying to fit their Rails mindset into a Go codebase. It takes a while, but it'll click eventually if you keep working at it.

Anyway, that's our basic API structure. The only way this really gets a lot more complicated is if you have dozens of endpoints.

For complex services with async workers and such, we'll just add more directories in the /internal path to satisfy whatever we need to do, keeping it domain-oriented as much as we can.

Hope that helps.

## Avoiding SQL Injection
https://go.dev/doc/database/sql-injection
https://stackoverflow.com/questions/26345318/how-can-i-prevent-sql-injection-attacks-in-go-while-using-database-sql

## New CSS
https://mxb.dev/blog/old-dogs-new-css-tricks/
