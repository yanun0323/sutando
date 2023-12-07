<a href="."><img height="200" src="https://github.com/yanun0323/asset/blob/main/sutando.png?raw=true"></a>
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
- Sort
- Limit
- Skip

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
        ClientOptionsHandler func(opts *options.ClientOptions) {
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
        ClientOptionsHandler func(opts *options.ClientOptions) {
            // do something...
        },
    })
```

- Model Declaration
```go
    // Supported
    type Element struct {
        FirstName string                                // 'firstName' as mongo db field key 
        lastName string                                 // 'lastName' as mongo db field key
        Nickname bool               `bson:"nick_name"`  // using `bson:"xxx"` tag to assign field key to 'xxx'
        Healthy bool                `bson:"-"`          // using `bson:"-"` tag to ignore this field
        Children []string           `bson:",omitempty"` // using `bson:",omitempty"` tag to ignore this field when it's empty
        CareerPlan CustomStruct                         // 'careerPlan' as mongo db field key  works
        Hobbies map[string]string
        Live time.Time
        Salary decimal.Decimal
    }
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
|1.3.3| - Added `Sort` `Limit` `Skip` into `Find`
|1.3.2| - Added deprecated comment for `DB`
|1.3.1| - Renamed `OptionHandler` to `ClientOptionsHandler` <br> - Renamed `SetupOption` to `SetupClientOptions`
|1.3.0| - Added `Execute Chain` <br> - Fixed error when input only one slice in insert function <br> - Fixed error when input only one param/slice in In/NotIn function <br> - Fixed `bson` `omitempty` supported <br> - Fixed embed structure lowercase Name issue <br> - Fixed map structure value lowercase Name issue <br> - Fixed array structure value lowercase Name issue <br> - Plan to remove db execute function in version 1.4.X
|1.2.1| - Support `mongodb-srv` <br> - Fixed `Conn` `ClientOptionsHandler` nill pointer issue
|1.2.0| - Added `ClientOptionsHandler` into `Conn` Interface
|1.1.2| - Fixed testing structure tag issue <br> - Fixed error wrapping issue
|1.1.1| - Added `Disconnect` function
|1.0.4| - Fixed some testing mistakes
|1.0.3| - Added `NewDBFromMongo` function
|1.0.2| - Added MIT License <br> - Removed Makefile
|1.0.1| - Fixed some testing mistakes
|1.0.0| - Release

## License

[MIT](https://github.com/yanun0323/sutando/blob/master/LICENSE)
