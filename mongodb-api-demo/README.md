## Set up the project
```
>> go mod init mongodb-api-demo
>> go get github.com/gin-gonic/gin
>> go get go.mongodb.org/mongo-driver
>> go mod tidy
```

## Run the project
```
>> go run main.go
```

## Knowledge
### MongoDB bson
`D` is an ordered representation of a BSON document. This type should be used when the order of the elements matters, such as MongoDB command documents. If the order of the elements does not matter, an `M` should be used instead.  
Example:
```
bson.D{{"foo", "bar"}, {"hello", "world"}, {"pi", 3.14159}}
```
`E` represents a BSON element for a `D`. It is usually used inside a `D`.

### Function init
In Go, the init function is a special function that is run before the main function of a package. It is used to initialize variables and do other setup tasks that need to be done before the main function is run. The init function is called automatically when a package is imported, so you don't have to call it yourself.

```
package main

import "fmt"

var message string

func init() {
	message = "Hello, World!"
}

func main() {
	fmt.Println(message)
}

```

### Connect to MongoDB
```
mongoConnection := options.Client().ApplyURI("mongodb+srv://<username>:<password>@<hostname>/?retryWrites=true&w=majority")
mongoClient, err := mongo.Connect(ctx, mongoConnection)
```

### Get sensitive variables from environment
In CMD:
```
>> export mongodb_username='username'
>> echo ${mongodb_username}
```

In Golang:
```
username := os.Getenv("mongodb_username")
```
