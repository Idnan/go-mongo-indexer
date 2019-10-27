# go-mongo-indexer

> CLI utility to manage mongodb collection indexes

[![Blog URL](https://img.shields.io/badge/Author-blog-green.svg?style=flat-square)](https://adnanahmed.info)
[![Build Status](https://travis-ci.org/idnan/go-mongo-indexer.svg?branch=master)](https://travis-ci.org/idnan/go-mongo-indexer)
[![Go Report Card](https://goreportcard.com/badge/github.com/gohugoio/hugo)](https://goreportcard.com/report/github.com/gohugoio/hugo)

## Usage

```shell
indexer --config <index-config-file> 
        --uri <mongodb-connection-uri>
        --database <database name>
        --apply
```

Details of options is listed below

| **Option** | **Required?** | **Description**                                                                                              |
|------------|---------------|--------------------------------------------------------------------------------------------------------------|
| `config`   | Yes           | Path to [indexes configuration file](#config-format)                                                         |
| `uri`      | Yes           | MongoDB connection string e.g. `mongodb://127.0.0.1:27017`                                                   |
| `database` | Yes           | Database name                                                                                                |
| `apply`    | No            | Whether to apply the indexes on collections or not. If not given, it will show the plan that will be applied |


## Config Format

The configuration file is just a simple json file containing the indexes to be applied. This file is an array of objects. Where each object has details like collection name, cap size and indexes for this specific collection.
```javascript
[
    {
        "collection": "order",     // name of collection
        "cap": null,               // Number of bytes 
        "index": [                 // Array of index details
            {"cartId": 1},         // An ascending order index
            {"status": -1},        // Descending order index
            {"orderId": 1},
            {"groupId": 1},
            {"currency": 1},
            {"createdAt": -1},
            {"orderNumber": 1, "type": 1},  // Composite index on orderNumber and type
            {"type": -1, "paymentStatus": -1, "payment.paymentMethod": -1}
        ]
    },
    {
        "collection": "collection_name",
        "cap": null,
        "index": [
            {"userId": -1},
            {"username": 1},
            {"orderId": 1, "_unique": 1},                       // creates a `unique index`
            {"createdAt": -1, "_expireAfterSeconds": 500}       // creates a `expires index` that will delete document after given number of seconds 
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
indexer --config "/path/to/xyz.json" --uri "mongodb://127.0.0.1:27017/database_name" --database "database_name"
```

<p align="center">
        <img src="https://i.imgur.com/3yj4gMh.png" height="400px"/>
</p>

> Apply the index changes
```shell
$ indexer --config "/path/to/xyz.json" --uri "mongodb://127.0.0.1:27017/database_name"  --database "database_name" --apply
```

## Todo
* [ ] Write tests
* [x] Collection capping
* [x] Support for `_unique` and `_expireAfterSeconds` indexes

## Contributing

Anyone is welcome to contribute, however, if you decide to get involved, please take a moment to review the guidelines:

* Only one feature or change per pull request
* Write meaningful commit messages
* Follow the existing coding standards

## License
MIT © [Adnan Ahmed](https://github.com/idnan)
