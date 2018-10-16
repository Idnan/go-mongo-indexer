# go-mongo-indexer

> CLI utility to manage mongodb collection indexes

[![Go Report Card](https://goreportcard.com/badge/github.com/idnan/go-mongo-indexer)](https://goreportcard.com/report/github.com/idnan/go-mongo-indexer)
[![Blog URL](https://img.shields.io/badge/Author-blog-green.svg?style=flat-square)](https://adnanahmed.info)
[![Build Status](https://travis-ci.org/idnan/go-mongo-indexer.svg?branch=master)](https://travis-ci.org/idnan/go-mongo-indexer)

## Usage

```shell
indexer --config <index-config-file> 
        --uri <mongodb-connection-uri>
        --apply
```

Details of options is listed below

| **Option** | **Required?** | **Description** |
|------------|--------|-------|
| `config` | Yes | Path to [indexes configuration file](#config-format) |
| `uri`    | Yes | MongoDB connection string e.g. `mongodb://db1.example.net:27017` |
| `apply`  | No  | Whether to apply the indexes on collections or not. If not given, it will show the plan that will be applied |


## Config Format

The configuration file is just a simple json file containing the indexes to be applied. This file is an array of objects. Where each object has details like collection name, cap size and indexes for this specific collection.
```javascript
[
    {
        "collection": "order",     // name of collection
        "cap": null,               // Number of bytes 
        "index": [                 // Array of index details
            ["cartId"],            // An ascending order index
            ["-status"],           // Descending order index
            ["orderId"],
            ["groupId"],
            ["currency"]
            ["-createdAt"],
            ["orderNumber", "type"],  // Composite index on orderNUmber and type
            ["-type", "-paymentStatus", "-payment.paymentMethod"]
        ]
    },
    {
        "collection": "collection_name",
        "cap": null,
        "index": [
            ["-userId"],
            ["username"],
            ["-createdAt", "-user.email"]
        ]
    }
    ....
    ....
    ....
]
```

**Note** `cap`ping is still progress not yet supported

## Examples

> See list of index changes before applying

```shell
indexer --config "/path/to/xyz.json" --uri "mongodb://127.0.0.1:27017/database_name"
```

![plan](https://i.imgur.com/3yj4gMh.png)

> Apply the index changes
```shell
$ indexer --config "/path/to/xyz.json" --uri "mongodb://127.0.0.1:27017/database_name" --apply
```

## Todo
* [ ] Write tests
* [ ] Collection capping
* [ ] Support for `unique` and `expireAt` indexes

## Contributing

Anyone is welcome to contribute, however, if you decide to get involved, please take a moment to review the guidelines:

* Only one feature or change per pull request
* Write meaningful commit messages
* Follow the existing coding standards

## License
MIT © [Adnan Ahmed](https://github.com/idnan)
