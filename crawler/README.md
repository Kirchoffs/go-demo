# Crawler

## Create the Project
```
>> mkdir crawler && cd crawler
>> go mod init crawler
```

## Add Dependencies
```
>> go get golang.org/x/net/html/charset
>> go get golang.org/x/text/encoding
```

```
>> go list -m -versions github.com/PuerkitoBio/goquery
>> go list -m all
```

## Run the Project
```
>> go build main.go
>> ./main
```

```
>> go run main.go > thepaper.html
```

## Test
```
>> go test ./...
>> go test ./tests  -v -run TestChromeDP
```

## Code Details
### strings & bytes
```go
import strings

body := io.ReadAll(resp.Body)

numLinks := strings.Count(string(body), "<a")
exist := strings.Contains(string(body), "<a")
```

```go
import strings

body := io.ReadAll(resp.Body)

numLinks := bytes.Count(body, []byte("<a"))
exist := bytes.Contains(body, "<a")
```

## Golang Details
### Format Print
- __%v__ is a placeholder that represents the value of an operand. It is used to interpolate the value of a variable into a string. The behavior of %v depends on the type of the operand being printed.

### Context with Web
New methods related to context and web development:

- __Context()__ returns the context.Context associated with the request.

- __WithContext(ctx)__ takes in a context.Context and returns a new http.Request with the old request’s state combined with the supplied context.Context.

```go
func Middleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        // wrap the context with stuff
        req = req.WithContext(ctx)
        handler.ServeHTTP(rw, req)
    })
}
```

```go
func handler(rw http.ResponseWriter, req *http.Request) {
    ctx := req.Context()
    err := req.ParseForm()
    if err != nil {
        rw.WriteHeader(http.StatusInternalServerError)
        rw.Write([]byte(err.Error()))
        return
    }
    data := req.FormValue("data")
    result, err := process(ctx, data)
    if err != nil {
        rw.WriteHeader(http.StatusInternalServerError)
        rw.Write([]byte(err.Error()))
        return
    }
    rw.Write([]byte(result))
}
```

Use `NewRequestWithContext` from net/http package:
```go
type ServiceCaller struct {
    client *http.Client
}

func (sc ServiceCaller) callAnotherService(ctx context.Context, data string) (string, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://target.com?data=" + data, nil)
    if err != nil {
        return "", err
    }

    resp, err := sc.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("Unexpected status code %d", resp.StatusCode)
    }
    
    // do the rest of the stuff to process the response
    id, err := processResponse(resp.Body)
    return id, err
}
```