# Scratchpad

- the [repo]('https://github.com/ardanlabs/service') for the ultimate go: service with k8 course

## Module 2
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

## Module 3
- A sidecar is kinda like a pod that runs alongside another pod. So if your service needs a metrics pod, but only when the sevice pod is running, you can run the metrics pod as a sidecar.
- Zipkin is a tracing sidecar

## Module 4
- not sure what I think about subpackages not importing each other. On one hand, I see that, but it feels like the level of organization he is suggesting is level 10/10 and I don't know how you could come up with a working design that would work for every project.
- Interesting to pick directory names that keep things in order from top to bottom
- I absolutely HATE having dockerfiles live so far away from the go they are trying to build
- Being able to understand something is more important than it being easy to do
- validate signal to noise ratio on logs, ideally the logs in dev are the same as the logs in prod
- like the idea of everything having a default setting that makes things work at least in dev