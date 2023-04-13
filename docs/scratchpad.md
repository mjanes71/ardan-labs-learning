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