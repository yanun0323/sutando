# Sutando
A go package encapsulate mongo-driver.


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
    query := db.Collection("Collection").Find().Equal("Name", "sutando").Greater("Number", 300).First()
    err := db.ExecFind(ctx, query, &result)
```
#### Create
```go
    insert := db.Collection("Collection").Insert(&obj)
    result, _, err := db.ExecInsert(ctx, insert)

    insertMany := db.Collection("Collection").Insert(&obj1, &obj2, &obj3)
    _, resultMant, err := db.ExecInsert(ctx, insertMany)
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
    mongoClient := db.Client()
    mongoDB := db.DB()
``` 