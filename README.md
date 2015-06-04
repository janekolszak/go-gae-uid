# go-gae-uid
Used the code from: https://cloud.google.com/appengine/articles/sharding_counters
And this library https://github.com/speps/go-hashids for generating hash uids.

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


