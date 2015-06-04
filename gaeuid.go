package gaeuid

import (
    "appengine"
    "appengine/datastore"
    "appengine/memcache"
    "fmt"
    "github.com/speps/go-hashids"
    "math/rand"
)

const (
    defaultShards = 25
    configKind    = "_GAEUIDShardConfig"
    shardKind     = "_GAEUIDShard"
)

type Generator struct {
    name   string
    hashID *hashids.HashID
}

// Internal structures
type counterConfig struct {
    Shards int
}

type shard struct {
    Name  string
    Count int64
}

func NewGenerator(name string, salt string, minLength int) *Generator {
    hd := hashids.NewData()
    hd.Salt = salt
    hd.MinLength = minLength

    return &Generator{
        name:   name,
        hashID: hashids.NewWithData(hd),
    }
}

func (gen *Generator) NewID(c appengine.Context) (string, error) {
    value, err := gen.count(c)
    if err != nil {
        c.Errorf("GAEUID, Error in NewID: %s", err)
        return "", err
    }

    err = gen.increment(c)
    if err != nil {
        c.Errorf("GAEUID, Error in NewID: %s", err)
        return "", err
    }

    return gen.hashID.EncodeInt64([]int64{value})
}

// IncreaseShards increases the number of shards for the named counter to n.
// It will never decrease the number of shards.
func (gen *Generator) IncreaseShards(c appengine.Context, n int) error {
    ckey := datastore.NewKey(c, configKind, gen.name, 0, nil)
    return datastore.RunInTransaction(c, func(c appengine.Context) error {
        var cfg counterConfig
        mod := false
        err := datastore.Get(c, ckey, &cfg)
        if err == datastore.ErrNoSuchEntity {
            cfg.Shards = defaultShards
            mod = true
        } else if err != nil {
            c.Errorf("GAEUID, Error in IncreaseShards: %s", err)
            return err
        }
        if cfg.Shards < n {
            cfg.Shards = n
            mod = true
        }
        if mod {
            _, err = datastore.Put(c, ckey, &cfg)
        }
        c.Errorf("GAEUID, Error in IncreaseShards: %s", err)
        return err
    }, nil)
}

func memcacheKey(name string) string {
    return shardKind + ":" + name
}

// Count retrieves the value of the named counter.
func (gen *Generator) count(c appengine.Context) (int64, error) {
    var total int64
    total = 0

    mkey := memcacheKey(gen.name)
    if _, err := memcache.JSON.Get(c, mkey, &total); err == nil {
        return total, nil
    }
    q := datastore.NewQuery(shardKind).Filter("Name =", gen.name)
    for t := q.Run(c); ; {
        var s shard
        _, err := t.Next(&s)
        if err == datastore.Done {
            break
        }
        if err != nil {
            return total, err
        }
        total += s.Count
    }
    memcache.JSON.Set(c, &memcache.Item{
        Key:        mkey,
        Object:     &total,
        Expiration: 60,
    })
    return total, nil
}

// Increment increments the named counter.
func (gen *Generator) increment(c appengine.Context) error {
    // Get counter config.
    var cfg counterConfig
    ckey := datastore.NewKey(c, configKind, gen.name, 0, nil)

    err := datastore.RunInTransaction(c, func(c appengine.Context) error {
        err := datastore.Get(c, ckey, &cfg)
        if err == datastore.ErrNoSuchEntity {
            cfg.Shards = defaultShards
            _, err = datastore.Put(c, ckey, &cfg)
        }
        return err
    }, nil)

    if err != nil {
        return err
    }

    // Pick one shard at random and increment
    var s shard
    err = datastore.RunInTransaction(c, func(c appengine.Context) error {
        shardName := fmt.Sprintf("%s-shard%d", gen.name, rand.Intn(cfg.Shards))
        key := datastore.NewKey(c, shardKind, shardName, 0, nil)
        err := datastore.Get(c, key, &s)
        // A missing entity and a present entity will both work.
        if err != nil && err != datastore.ErrNoSuchEntity {
            return err
        }
        s.Name = gen.name
        s.Count++
        _, err = datastore.Put(c, key, &s)
        return err
    }, nil)
    if err != nil {
        return err
    }
    memcache.IncrementExisting(c, memcacheKey(gen.name), 1)
    return nil
}
