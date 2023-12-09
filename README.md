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
- Contain
- In
- NotIn
- Sort
- Limit
- Skip
- Count
- Regex

## Guide

### Installation

```shell
$ go get -u github.com/yanun0323/sutando@latest
```

### Example

#### Connect To MongoDB

- Create a new connection
```go
    // connect using host and port.
    db, err := sutando.NewDB(ctx, sutando.Conn{
    	Username:  "example",
    	Password:  "example",
    	Host:      "example",
    	Port:      27017,
    	DB:        "example",
    	AdminAuth: true,
    	Pem:       "",
    	ClientOptionsHandler: func(opts *options.ClientOptions) {
    		opts.SetConnectTimeout(5 * time.Second)
    		opts.SetTimeout(15 * time.Second)
    	},
    })

    // connect using SRV url.
    db, err := sutando.NewDB(ctx, sutando.ConnSrv{
    	Username:  "example",
    	Password:  "example",
    	Host:      "example.mongo.net",
    	DB:        "example",
    	AdminAuth: true,
    	Pem:       "",
    	ClientOptionsHandler: func(opts *options.ClientOptions) {
    		opts.SetConnectTimeout(5 * time.Second)
    		opts.SetTimeout(15 * time.Second)
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

#### Drop
```go
    err := db.Collection("Collection").Drop(ctx)
```

#### Scalar
```go
    // Count
    count, err := db.Collection("Collection").Find().Equal("Name", "sutando").Greater("Number", 300).Count(ctx, "_index_id_")
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
    client := db.RawClient()
    database := db.RawDatabase()
``` 

## Changelog

|Version|Description
|:-:|:-
|1.4.0| - Remove all db execute functions <br> - Removed `Query` <br> - Removed method `Bitwise` <br> - Added `ConnSrv` connection structure <br> - Added method `Drop` into `Collection()` method chain <br> - Added comment for all methods <br> - Added `option` parameter into `Regex` method <br> - Rewrite the structure fo filters <br> - Renamed `GetDriver` to `RawClient` <br> - Renamed `GetDriverDB` to `RawDatabase` <br> - Fixed after invoking `Find`, didn't call `defer cursor.Close()`
|1.3.7| - Added `Scalar` <br> - Moved method `Count` from `Query` to `Scalar`
|1.3.6| - Added method `Regex` into `Update` `Find` `Delete` `Query`
|1.3.5| - Fixed `Find` no document mismatch error
|1.3.4| - Added method `Count` into `Query`
|1.3.3| - Added methods `Sort` `Limit` `Skip` into `Find`
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
