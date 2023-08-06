<a href="."><img height="200" src="https://raw.githubusercontent.com/yanun0323/asset/main/sutando.png?token=GHSAT0AAAAAACFQBOSGHZQSMQV6U2DHKREYZGM2GUQ"></a>

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
    resultOne := struct{}
    err := db.Collection("Collection").Find().Equal("Name", "sutando").Greater("Number", 300).First().Exec(ctx, &resultOne)

    resultMany := []struct{}
    err := db.Collection("Collection").Find().Equal("Name", "sutando").Greater("Number", 300).Exec(ctx, &resultMany)
```
#### Create
```go
    resultOne, _, err := db.Collection("Collection").Insert(&obj).Exec(ctx)

    _, resultMany, err := db.Collection("Collection").Insert(&obj1, &obj2, &obj3).Exec(ctx)
```
    
#### Update with Model (Will update all fields including empty fields)
```go
    resultOne, err := db.Collection("Collection").UpdateWith(&data).Equal("Field", "sutando").First().Exec(su.ctx, false)

    resultMany, err := db.Collection("Collection").UpdateWith(&data).Equal("Field", "sutando").Exec(su.ctx, false)
```
#### Update with Set
```go
    resultOne, err := db.Collection("Collection").Update().Equal("Field", "sutando").First().Set("Field", "hello").Exec(su.ctx, false)

    resultMany, err := db.Collection("Collection").Update().Equal("Field", "sutando").Set("Field", "hello").Exec(su.ctx, false)
```
#### Delete
```go
    resultOne, err := db.Collection("Collection").Delete().Equal("Field", "sutando").First().Exec(su.ctx)

    resultMany, err := db.Collection("Collection").Delete().Equal("Field", "sutando").Exec(su.ctx)
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