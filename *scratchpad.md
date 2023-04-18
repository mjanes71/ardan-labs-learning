# Scratchpad

- the [repo]('https://github.com/ardanlabs/service') for the ultimate go: service with k8 course

## Module 1
- name your mod file starting with your repo
- vs code can be confused if you don't have a go mod (might open it in cantral gopath mode). Now that we have mods, the gopath goes to the level of the go.mod
- cmd + shift + p or reload the window to reload the language server
- At a systems level, y'all need to center on "what is a project". For this teacher, a project represents a repo.
- One project can contain multiple binaries. I think that's gonna be why we need multiple go mods per project/repo.
- A project will have conventions that everyone needs to conform to, even if they don't love it.
- Apparently go has a backwards compatability promise (1.16 is compatable with everything before 1.16)
- you can access all your go env variables like:
```go
go env
```
- the `GOMODCACHE` is the path at which all dependency source code is stored by go
- gopls is the language server (the go team wrote) that vs code is using to give you intellisense and other stuff in vscode. it's also storing some go caches (like the modcache)
- to add stuff to the modcache locally, use `go mod tidy` which will run some go gets for you as it looks at your project to see what imports you're trying to use in your code
- the `GOPROXY` tells go where to go to look for source code. Proxy server = module mirror. Serves as a proxy for github, gitlab, bitbucket, etc. So when you go mod tidy, you're getting a list of every version the proxy server has. If the proxy server doesn't have it, it looks at the source host (like GH) to get that code. By default, it will choose the latest and greatest.
- the go sum stores hash codes that allow us to validate that the source code we got for a module is what we would expect it to be
- `GONOPROXY` is a variable that you can set to tell go directly to another server to get source code for certain imports like "github.com/pgdevelopers"
- There are other proxy servers, like Athens, that you might choose to use instead (you can run and install it on your own network). If you used that, you would just change your GOPROXY
- ChecksumDB is run by the go team. If you're writing a hash to the go sum for the first time, it checks the checksum db to compare the hash codes and make sure they match. If its already in go sum, it just compares against what it already has. Makes sure that no matter what source you end up pulling from, its always gonna match or you're gonna get an error.
- When you go mod vendor, go stops looking at the module cache and starts looking at the vendor dir. Big benefits here are you can see and even hack on that vendored code a little bit. Building also gets easier because you don't need the internet to build (cuz you don't have to download anything). Once you have a vendor dir at the same level as the mod, you need to run a go mod vendor after a go mod tidy to keep things in sync
- If you ever see an "incompatable" on an import, you might have to specify the import with a version like `github.com/mykewlpkg/v4` or something like that. If you just do `github.com/mykewlpkg` its gonna grab v1
- Go, when choosing version of downstream dependencies, may not always choose the latest and greatest versions of a downstream dependency. This could cause a problem if another dependency requires a different version. The GO algo will choose the lowest version that will work for everything
- If you want to use the latest and greatest versions of things all the time, you can do `go get -u -t -v ./...` This walks the project and gets the latest greatest of every dependency. Then you tidy and vendor.
- If you pass a version with a dependency, you can maintain two different versions of that source code in your mod
- A sidecar is kinda like a pod that runs alongside another pod. So if your service needs a metrics pod, but only when the sevice pod is running, you can run the metrics pod as a sidecar.
- Zipkin is a tracing sidecar
- Deploy first mentality = getting just enough code to get a service that runs
- Understand where you fit in the knowledge spectrum of your project  (if you're above average, write and design to the middle, if stuff is too clever, other people with less knowledge can't maintain it)
- you always want to have no more go processors running than you have machines/cores
- ?? What do y'all think about putting stuff in a run function? Can someone explain this to me better than Bill? 

## Module 2
- think about what hat you're wearing when you write code. devops vs dev. what if some day you aren't both?
- sidecars: separate metrics or other intellegence from the main service (devops could write sidecars and be in charge of where that stuff goes to (like honeycomb))
- kubectl apply is a command where k8 cluster is taking your instruction and figuring out how to make it work (asynchronous, won't happen immediately)
- deployment defines what the pod needs to look like
- service defines the networking side of things
- some go dependencies might use c libraries by default, but if you disable it, they have a go native workaround. bill likes to disable cgo in his test commands to make sure the docker build (with cgo disabled) is gonna run the same way things are locally
- i learned from [this article](https://levelup.gitconnected.com/a-better-way-than-ldflags-to-add-a-build-version-to-your-go-binaries-2258ce419d2d) that you can use go embed and go generate instead of ldflags
- ??? thoughts on addgroup and adduser in dockerfiles? being explicit about access?
- for kustomization setup, the kustomization.yaml in base points to the yaml in base. the kustomization.yaml in dev/sales points to the kustomization.yaml in base AND the patch its trying to patch in
- the patch metadata needs to match the base yaml matadata so kustomize understands what yaml its trying to patch
- the go runtime is not k8 aware, so if you are cpu restricting in your yaml spec, you need to also use the maxprocs library to help go be aware of the k8 configs
- Bills rules about configs:
1. only main.go should be referencing configs
2. all configs should have a dev default with few exceptions (like keys)
3. the config service should support --help
- ??? Thoughts on using the ardan labs conf pkg?
- the logfmt tool is some next level $h!t. i wonder when he decides its time to do this stuff. or did he just steal it from someone else.
- i feel like i'd have to have 800 years of go experience to get on Bill's level
- a mux is a piece of code / router that accepts requests, and if it has a matching route, it will route the request to the right handler. watchout on user the defaultmuxlibrary, its a security risk
- think of goroutines in a parent-child relationship. parent goroutines should not terminate until its sure the child routines have terminated. one approach is to have a pkg with a waitgroup that is called by the top level parent to monitor shutdown of children

## Module 3
- a goroutine is just a function with go in front of it. this tells go that you want to run that bit of code asynchronously with any other goroutines in the code. ie:
```go
// Prints numbers from 1-3 along with the passed string
func foo(s string) {
    for i := 1; i <= 3; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Println(s, ": ", i)
    }
}

func main() {
    
    // Starting two goroutines
    go foo("1st goroutine")
    go foo("2nd goroutine")

    // Wait for goroutines to finish before main goroutine ends
    time.Sleep(time.Second)
    fmt.Println("Main goroutine finished")
}
```
- all go routines are running within the context of the main func. if you shutdown main without accounting for any running goroutines, you risk orphaning a routine
- load-shedding using http package is how you make sure stuff doesn't get orphaned/shutdown before its time
- !!! DEFININTELY gonna steal the convention of adding a comments section at the top under imports where todos can go instead of having them spread out all over the doc where i'm likely to miss them. especially if something us under active dev. could even have the linter check to see if there are things left in this and alert you probably.
- Bill thinks the backlog is where things go to die. I rather enjoy grooming the backlog.
- synchronizing access to shared state is like waiting in a line to place an order for coffee. the person at the front of the line as access to the state currently. atomic instructions and mutexes
- when you get to the front of the line, you're in orchestration mode
- channels are used to help a goroutine communicate back to the main func that its still running (this is the arrow syntax)
- use a channel to monitor signals (singal if ther is an error with your api, or a shutdown for example)
- ??? pattern of passing config object to function call, wonder which people like better
```go 

// option 1
cfgMux := handlers.APIMuxConfig{
		Shutdown: shutdown,
		Log: log,
	}
apiMux := handlers.APIMux(cfgMux)


//option 2
apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown: shutdown,
		Log: log,
	})

// applies to cockroach svc like:
// option 1
cnxInfo := cockroachconnect.CnxStrConfig{
			Driver: "postgres",
			User: os.Getenv("CR_USER"),
			Password: os.Getenv("CR_PASSWORD"),
			Host: os.Getenv("CR_HOST"),
			Port: "26257",
			DbName: os.Getenv("CR_DB"),
			SSLMode: "verify-full",
			SSLRootCert: os.Getenv("CR_ROOT_CERT"),
			AppName: "schema-migration",
		}
dbConnection, err = cockroachconnect.New(false, &cnxInfo)
if err != nil {
    log.Fatalf("Error connecting to database: %q", err)
}

// or option 2
dbConnection, err = cockroachconnect.New(false, cockroachconnect.CnxStrConfig{
			Driver: "postgres",
			User: os.Getenv("CR_USER"),
			Password: os.Getenv("CR_PASSWORD"),
			Host: os.Getenv("CR_HOST"),
			Port: "26257",
			DbName: os.Getenv("CR_DB"),
			SSLMode: "verify-full",
			SSLRootCert: os.Getenv("CR_ROOT_CERT"),
			AppName: "schema-migration",
		})
if err != nil {
    log.Fatalf("Error connecting to database: %q", err)
}

```

- ctx should be the first arg in any function that has i/o or a role in i/o
- ?? wonder if I could have used his "write a mux without writing a mix" fix to extend the the datadog mocks
- middleware functions are designed to take something and wrap it with something else. ie: take some existing function and wrap it with logging, or otel, or something else
- ??? wondering how common it is to reference go types you created in other files in the same pkg like this:
```go
package web

// Middleware is a function designed to run some code before and/or after
// another Handler. It is designed to remove boilerplate or other concerns not
// direct to any given Handler.
type Middleware func(Handler) Handler
```
- every directory is a pkg, every pkg is an api. one file named after the package name so its clear what the pkg does
- utils, common, helpers, models (you end up having a fragile codebase sometimes from things like this). just be careful how much stuff is importing that pkg, more imports = more fragility
- ???trying to understand closures (when an anon function closes over a value)
```go
// Logger writes information about the request to the logs.
func Logger(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			log.Infow("request started", "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Infow("request completed", "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			return err
		}

		return h
	}

	return m
}
```

- i think working with Bill would be awesome and so frustrating. it feels like he has exceptions to every rule, but his rules are so hard and numerous
- ??? what strategies is anyone employing to get to Bill's level of go understanding/expertise? It's way too much to process this entire course and implement it all in your own day-to-day. I'm inclined to take one thing, like logging, and do a deep dive on it, and then develop my own opinions on it and start implementing them as I write go.
- errors should only be handled once
- once an error is handled, it should be done
- errors should probably be handled at the highest layers (business, or app in this project)
- error is an interface. when we say `if err != nil` what we're asking is, has a concrete value been stored in the error interface you returned? If something has been stored there, it indicates an error value that needs to be handled.
- Bill rocking back and forth from side to side at 2x makes me want to vom
- the lower layers should not shut down/panic
- creating structs and using pointer receivers to build out their functionality protects from fraud. Bill explains it decently well.
- when you've been listening to bill on 2x and take it down to 1x, it feels like you're living in slowmo