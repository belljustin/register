# Compiletime and Runtime drivers

This post explores adding and building drivers both at compiletime and runtime.

# Go Register

When building web apps with Go, one of the first packages you probably encounter is `database/sql`.
This package must be used in conjunction with a database driver like `lib/pq` for postgres.

The documentation for `lib/pq` provides a handy example for getting started.

```go
import (
	"database/sql"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	age := 21
	rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
	…
}
```

An interesting point about this snippet is the import with a leading `_` - an "anonymous import".
Why would we want to import a library if we don't plan on calling it?
Well, the imported package still runs initialization meaning all variable declarations and `init()` functions are evaluated.
If we inspect the `init()` of [pg/conn.go](https://github.com/lib/pq/blob/8446d16b8935fdf2b5c0fe333538ac395e3e1e4b/conn.go#L57-L59), we see it interacts with a `Register(...)` method in `database/sql`.

```go
func init() {
	sql.Register("postgres", &Driver{})
}
```

This is handy because it provides our common interface of `database/sql` with a driver for postgres.
If you check out the source for `sql.Register(...)`, you see it just stores that driver in map keyed by the provided string.

```go
var (
	driversMu sync.RWMutex // mutex protects the drivers map during concurrent access
	drivers   = make(map[string]driver.Driver)
)

func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("sql: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	drivers[name] = driver
}
```

# Writing our own Drivers

Of course we can copy this pattern to create our own interfaces and drivers that magically wire themselves up.

Here's a very simple interface to a foreign exchange (forex) service.

```go
type Currency string

const (
	USD Currency = "USD"
	CAD Currency = "CAD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
)

// ForexService is an interface that provides methods for querying foreign exchange rates.
type ForexService interface {
	// GetRate returns the exchange rate of a currency pair in basis points or an error if one occurred.
	GetRate(c1, c2 Currency) (int, error)
}
```

I've implemented drivers for this interface using both [fawazahmed0/currency-api](https://github.com/fawazahmed0/currency-api) and [freeforexapi](https://freeforexapi.com/).
Just like the database example in the standard lib, you can anonymously import one of the drivers and use the ForexService.

```go
package main

import (
	"fmt"

	"github.com/belljustin/register"
	_ "github.com/belljustin/register/drivers/freeforex"
)

func main() {
	forexService := register.Open("freeforex")
	rate, _ := forexService.GetRate(register.USD,register.CAD)
	fmt.Println(rate)  // outputs the rate in bips
}
```

# Dynamic Loading

Building the above creates an executable with fixed behaviour at compile time.
However, the Go compiler also features a buildmode that creates a Go plugin.

```shell
go build -buildmode=plugin -o ./bin/freeforex.so ./drivers/freeforex/plugin
```

You can then load this plugin at runtime using the std lib.

```go
plugin.Open("./bin/freeforex.so")
```

This will run all the initialization but not the `main()` function of the opened package.
Without compiling or deploying the application, you can add new functionality.

Though, it's worth noting if you try to re-open the same plugin path, it will simply return the existing path.
Worse still, if you recompile the same package to another path and try to open that, it will result in an error.

```shell
panic: plugin.Open("plugins/freeforexV2"): plugin already loaded
```

At the time of writing, closing an existing plugin and reloading is not a supported feature.
It will likely continue to be unsupported because one would have to keep track of all the references to the plugin and clean them up in order to avoid seg faults.
So this is not a way to do hot-reloading or continuously deploy code without downtime.

Nonetheless, it is a cool way to load new functionality at runtime.
Perhaps more practically, it allows for smaller binaries because you can build the main component, and individual plugins seperately.
Then users can link only the individual plugins they want to use.
