# go-mongo-indexer
> CLI tool to manage mongo database indexes through json files

## Usage
```shell
$ index --config <json config file> --uri <mongo database uri>
```
Below is the description of all the accepted options
- `--apply` not required, to apply the index changes
- `--config` required, path to json file containing the each collection index
- `--uri` required, mongo database uri

## Examples

> See list of index changes before applying
```bash
$ index --config "xyz.json" --uri "mongodb://127.0.0.1:27017/database_name"

```

> Apply the index changes
```bash
$ index --config "xyz.json" --uri "mongodb://127.0.0.1:27017/database_name" --apply
``` 

## Todo
* [ ] write tests
* [ ] cap on collection
* [ ] support for `unique` and `expireAt` indexes

## License
MIT © [Adnan Ahmed](https://github.com/idnan)