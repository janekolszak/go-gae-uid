# go-gae-uid

The library generates unique, string IDs of a given minimal length.
As an optimization it uses sharding counters - code copied from a [GAE tutorial](https://cloud.google.com/appengine/articles/sharding_counters).

**Depends on**:
* [hashid library](https://github.com/speps/go-hashids).

### Setup
<pre>go get github.com/janekolszak/go-gae-uid</pre>

### Example
```go
package app

import (
    "fmt"
    "github.com/janekolszak/go-gae-uid"
)

func main() {
    gen = gaeuid.NewGenerator("Kind", "HASH'S SALT", 11 /*id length*/)
    id, err = gen.NewID(c)
    if err != nil {
        return err
    }
    fmt.Println("Unique id: ", id)
}
```


