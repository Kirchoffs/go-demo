# Notes

## Golang Tools
### Environment
```
>> go env
>> go env GOPATH
```

### Golang version
```
>> go version
```

### Gofmt
Applied format to the whole project
```
>> gofmt -s -w .
```

### Check the dependency version
```
>> go list -m -versions github.com/goccy/go-json

>> go list -m all | grep go-json
```

### Update all dependencies to the newest
```
>> go get -u -t -v ./...
```

## Golang Details
### Interface
The interface type that specifies zero methods is known as the empty interface:
```
interface{}
```
An empty interface may hold values of any type. (Every type implements at least zero methods.)

`any` is a new predeclared identifier and a type alias of `interface{}`.

## Git
In .git/hooks/pre-commit, add the following to change tabs to spaces:

```
#!/bin/sh
python3 go-space-format.py

SCRIPT_EXIT_STATUS=$?

if [ $SCRIPT_EXIT_STATUS -ne 0 ]; then
    echo "Python script failed, aborting commit."
    exit 1
fi

git add .

exit 0
```

Also don't forget to run `chmod +x .git/hooks/pre-commit` to make the script executable.
