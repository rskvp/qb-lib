# DB Bolt

## Introduction to BBolt

Bolt is a pure Go key/value store inspired by [Howard Chu's](https://twitter.com/hyc_symas)
[LMDB project](http://symas.com/mdb/).

The goal of the project is to provide a simple, fast, and reliable database for projects that don't
require a full database server such as Postgres or MySQL.

[BBolt](https://pkg.go.dev/go.etcd.io/bbolt/) supports fully serializable transactions, ACID semantics, and lock-free
MVCC with multiple readers and a single writer.
Bolt can be used for projects that want a simple data store without the need to add large
dependencies such as Postgres or MySQL.

Bolt is a single-level, zero-copy, B+tree data store. This means that Bolt is optimized for fast read access and does
not require recovery in the event of a system crash. Transactions which have not finished committing will simply be
rolled back in the event of a crash.

The design of Bolt is based on Howard Chu's LMDB database project.
Bolt currently works on Windows, Mac OS X, and Linux.

### Basics

There are only a few types in Bolt: DB, Bucket, Tx, and Cursor.
The DB is a collection of buckets and is represented by a single file on disk.
A bucket is a collection of unique keys that are associated with values.

Transactions provide either read-only or read-write access to the database.
Read-only transactions can retrieve key/value pairs and can use Cursors to iterate over the dataset sequentially.
Read-write transactions can create and delete buckets and can insert and remove keys.
Only one read-write transaction is allowed at a time.

## DB Bolt

GGBolt is a wrapper implementation that allow use BBolt with a simple implementation.

```
    // create and open DB
    db := qb_bolt.NewBoltDatabase(config)
    err = db.Open()

    // get/create collection
    coll, err := db.Collection("big-coll", true)

    // insert/update item
    item := &map[string]interface{}{
        "_key": "1",
        "name": "Mario",
        "age":  22,
    }
    err = coll.Upsert(item)
```

GGBolt implementation assume that document's key is named "\_key".

## Use GGBolt for Cache

```
filename := "./db/expiring.dat"
_ = qb_paths.Mkdir(filename)
config := qb_bolt.NewBoltConfig()
config.Name = filename

db := qb_bolt.NewBoltDatabase(config)
err := db.Open()
if nil!=err{
    panic(err)
}

defer db.Close()

coll, err := db.Collection("cache", true)
if nil!=err{
    panic(err)
}

// set collection as expirable
coll.EnableExpire(true)

// add an expirable item to collection
item := map[string]interface{}{
    "_key":    qbc.Rnd.Uuid(),
    "name":    "NAME " + qbc.Strings.Format("%s", i),
}
item[qb_bolt.FieldExpire] = time.Now().Add(5 * time.Second).Unix()
err = coll.Upsert(item)

if nil!=err{
    panic(err)
}

```

DB Bolt can be used as a temporary cache repository because has a feature that
automatically deletes data whose expiration time has reached.

To enable a collection to check items for timed expiration you need:

- Enable a collection to check for expired fields
- Add "\_expire" field to collection (use constant "qb_bolt.FieldExpire")

```
// set collection as expirable
coll.EnableExpire(true)

...

// add "_expire" field
item[qb_bolt.FieldExpire] = time.Now().Add(5 * time.Second).Unix()
```

`_expire` field must be a unix timestamp, an int64 value.

## When to Use DB Bolt

GGBolt is great when you need a full embeddable cross-platform key-value pair database.

Mongo, Arango Redis and so on are a fair way greater than GGBolt, but are not embeddable and fully cross-platform
like pure Go code is. GGBolt is just pure Go code.
