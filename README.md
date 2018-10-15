# go-mongo-indexer
> CLI tool to manage mongo database collection indexes through json files

## Usage
```shell
$ indexer --config <json config file> --uri <mongo database uri>
```
Below is the description of all the accepted options
- `--apply` not required, to apply the index changes
- `--config` required, path to json file containing the each collection index. This json file should be in this [format](https://gist.github.com/Idnan/71e0478985a3aefa88f6502f83b28056#file-go-mongo-indexer-format-json)
- `--uri` required, mongo database uri

## Json Format
The config json file that you will pass to `--config` flag to index the database collection should be in specific format. This file is an array of objects. Where each object has details like collection name, cap size and indexes for this specific collection.
```
[
    {
        "collection": "collection_name",
        "cap": null,
        "index": [
            ["cartId"],
            ["-status"],
            ["orderId"],
            ["groupId"],
            ["currency"]
            ["-createdAt"],
            ["orderNumber", "type"],
            ["-type", "-paymentStatus", "-payment.paymentMethod"],
        ]
    },
    {
        "collection": "collection_name",
        "cap": null,
        "index": [
            ["-userId"]
            ["username"],
            ["-createdAt", "-user.email"],
        ]
    }
    ....
    ....
    ....
]
```

- `Collection` is the collection name
- `cap` collection cap size in bytes (not yet supported)
- `index` array of indexes `-` sign before a column refers to the descending order and no sign before column name refers to the ascending order  

## Examples

> See list of index changes before applying
```shell
$ indexer --config "xyz.json" --uri "mongodb://127.0.0.1:27017/database_name"

```

> Apply the index changes
```shell
$ indexer --config "xyz.json" --uri "mongodb://127.0.0.1:27017/database_name" --apply
``` 

## Todo
* [ ] write tests
* [ ] cap on collection
* [ ] support for `unique` and `expireAt` indexes

## License
MIT © [Adnan Ahmed](https://github.com/idnan)