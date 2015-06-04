# go-gae-uid

The library generates unique, string IDs of a given minimal length.
As an optimization it uses sharding counters - code copied from a [GAE tutorial](https://cloud.google.com/appengine/articles/sharding_counters).

### Dependencies
* [hashid](https://github.com/speps/go-hashids)

### Setup
<pre>go get github.com/janekolszak/go-gae-uid</pre>

### Example
```go
package app

import (
    "appengine"
    "github.com/janekolszak/go-gae-uid"
)

gen := gaeuid.NewGenerator("Kind", "HASH'S SALT", 11 /*id length*/)

func main(w http.ResponseWriter, r *http.Request){
    c := appengine.NewContext(r)

    // Get an id
    id, err = gen.NewID(c)
    if err != nil {
        return
    }
    c.Infof("Unique id: ", id)

    // By default there are 25 counters
    // You can only increase this number, e.g.:
    n := 120
    gen.IncreaseShards(c, n)
}
```


