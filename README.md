# Sutando
A go package that encapsulate mongo operation.


## Install
*Required go 1.18 up*

```shell
$ go get -u github.com/yanun0323/sutando
```

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

## Example

#### Connect To MongoDB
```go
    db, err := sutando.NewDB(ctx, sutando.Conn{
        Username:  "example",
        Password:  "example",
        Host:      "localhost",
        Port:      27017,	// leave empty if there's port in host
        DB:        "example",
        AdminAuth: true,
        Pem:       "",		// optional
    })
```

#### Find
```go
    result := struct{}
    query := e.Collection("Collection").Find().Equal("Name", "sutando").Greater("Number", 300).First()
    err := e.ExecFind(ctx, query, &result)
```
#### Create
```go
    insert := e.Collection("Collection").Insert(&obj)
    result, _, err := e.ExecInsert(ctx, insert)

    insertMany := e.Collection("Collection").Insert(&obj1, &obj2, &obj3)
    _, resultMant, err := e.ExecInsert(ctx, insertMany)
```
    
#### Update with Model
```go
    update := su.db.Collection("Collection").UpdateWith(&data).Equal("Field", "sutando").First()
    result, err := su.db.ExecUpdate(su.ctx, updateOne, false)

    updateMany := su.db.Collection("Collection").UpdateWith(&data).Equal("Field", "sutando")
    result, err := su.db.ExecUpdate(su.ctx, updateMany, false)
```
#### Update with Set
```go
    update := su.db.Collection("Collection").Update().Equal("Field", "sutando").First().Set("Field", "hello")
    result, err := su.db.ExecUpdate(su.ctx, updateOne, false)

    updateMany := su.db.Collection("Collection").Update().Equal("Field", "sutando").Set("Field", "hello")
    result, err := su.db.ExecUpdate(su.ctx, updateMany, false)
```
#### Delete
```go
    delete := su.db.Collection("Collection").Delete().Equal("Field", "sutando").First()
    result, err := su.db.ExecDelete(su.ctx, delete)

    deleteMany := su.db.Collection("Collection").Delete().Equal("Field", "sutando")
    result, err := su.db.ExecDelete(su.ctx, deleteMany)
```

#### User mongo-driver instance
```go
    mongoClient := db.Client()
    mongoDB := db.DB()
``` 