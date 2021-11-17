2.1 Modules

    Go community stardard method for managing third party dependencies.

    An example module workflow:

        * in order for the go ecosystem to build it, a project must be initialized for the module system

            `go mod init <namespace i.e.: github.com/jnkroeker/service>`
            
            this creates a go.mod file at the root of the project; it contains the module name and the Go version

        * start to add new third-party dependencies (i.e. in main.go)

            import "github.com/aradanlabs/conf"

            use the package name with a function called New. 

                func main () { conf.New() }

            run `go mod tidy` 

                walks thru project validating that all source code needed is available on disk in module cache; $GOMODCACHE in `go env`

            all modules fetched will be added to go.mod

            go.sum file added to root next to go.mod

                contains the checksums of the zip files downloaded for a third-party dependency


2.2 Module Mirrors; `go mod tidy` behind the scenes

    GOPROXY of `go env` tells the Go tooling where to fetch third-party libraries

    Connect to private repos 
    
        set GONOPROXY (list of the domains not to reach out to the GOPROXY for), instead go direct

        Athens open source proxy server 


2.4 Vendoring; `go mod vendor`

    The idea is to not use the third-arty code that is sitting in the module cache, but to bring all dependencies into the project

    This way the project owns all the source code it depends on. 


