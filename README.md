# Go Register

When building web apps with Go, one of the first packages you probably will encounter is `database/sql`.
This package must be used in conjunction with a database driver like `lib/pq` for postgres.

The documentation for `lib/pq` provides a handy example for getting started.

```
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

```
func init() {
	sql.Register("postgres", &Driver{})
}
```

This is really handy because it provides our common interface of `database/sql` with a driver for postgres.
If you check out the source for `sql.Register(...)`, you see it just stores that driver in map keyed by the provided string.
Additionally, there's a mutex in there to protect against race conditions.

```
var (
	driversMu sync.RWMutex
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


