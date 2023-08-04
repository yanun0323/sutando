<img height="200" src="TITLE.PNG">

## Requirement
*Required go 1.18 up*

## Query Parameters

- Exists
- And
- Equal
- NotEqual
- Greater
- GreaterOrEqual
- Less
- LessOrEqual
- Bitwise
- Contain
- In
- NotIn

## Guide

### Installation

```shell
$ go get -u github.com/yanun0323/sutando
```

### Example

#### Connect To MongoDB

- Create a new connection
```go
    // Using Host and Port
    db, err := sutando.NewDB(ctx, sutando.Conn{
        Username:  "example",
        Password:  "example",
        Host:      "example",
        Port:      27017,	// leave blank if there's port in host
        DB:        "example",
        AdminAuth: true,
        Pem:       "",		// optional
        OptionHandler func(client *options.ClientOptions) {
            // do something...
        },
    })

    // Using SRV URL
    db, err := sutando.NewDB(ctx, sutando.Conn{
        Username:  "example",
        Password:  "example",
        Host:      "example.mongo.net",
        DB:        "example",
        AdminAuth: true,
        Srv:       true,
        OptionHandler func(client *options.ClientOptions) {
            // do something...
        },
    })
```

- Use an exist connection
```go
    var client *mongo.Client
    ...
    database := "example"
    db := sutando.NewDBFromMongo(ctx, client, database)

```

#### Disconnect
```go
    err := db.Disconnect(ctx)
```

#### Find
```go
    result := struct{}
    query := db.Collection("Collection").Find().Equal("Name", "sutando").Greater("Number", 300).First()
    err := db.ExecFind(ctx, query, &result)
```
#### Create
```go
    insert := db.Collection("Collection").Insert(&obj)
    result, _, err := db.ExecInsert(ctx, insert)

    insertMany := db.Collection("Collection").Insert(&obj1, &obj2, &obj3)
    _, resultMany, err := db.ExecInsert(ctx, insertMany)
```
    
#### Update with Model (Will update all fields include empty fields)
```go
    update := db.Collection("Collection").UpdateWith(&data).Equal("Field", "sutando").First()
    result, err := db.ExecUpdate(su.ctx, updateOne, false)

    updateMany := db.Collection("Collection").UpdateWith(&data).Equal("Field", "sutando")
    result, err := db.ExecUpdate(su.ctx, updateMany, false)
```
#### Update with Set
```go
    update := db.Collection("Collection").Update().Equal("Field", "sutando").First().Set("Field", "hello")
    result, err := db.ExecUpdate(su.ctx, updateOne, false)

    updateMany := db.Collection("Collection").Update().Equal("Field", "sutando").Set("Field", "hello")
    result, err := db.ExecUpdate(su.ctx, updateMany, false)
```
#### Delete
```go
    delete := db.Collection("Collection").Delete().Equal("Field", "sutando").First()
    result, err := db.ExecDelete(su.ctx, delete)

    deleteMany := db.Collection("Collection").Delete().Equal("Field", "sutando")
    result, err := db.ExecDelete(su.ctx, deleteMany)
```

#### Use original mongo-driver instance
```go
    client := db.GetDriver()
    db := db.GetDriverDB()
``` 

## Changelog

|Version|Description
|:-:|:-
|1.2.1| - Support `mongodb-srv` <br> - Fixed `Conn` `OptionHandler` nill pointer issue
|1.2.0| - Added `OptionHandler` into `Conn` Interface
|1.1.2| - Fixed testing structure tag issue <br> - Fixed error wrapping issue
|1.1.1| - Added `Disconnect` function
|1.0.4| - Fixed some testing mistakes
|1.0.3| - Added `NewDBFromMongo` function
|1.0.2| - Added MIT License <br> - Removed Makefile
|1.0.1| - Fixed some testing mistakes
|1.0.0| - Release

## License

[MIT](https://github.com/yanun0323/sutando/blob/master/LICENSE)