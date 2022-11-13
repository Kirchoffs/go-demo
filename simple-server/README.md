# Notes

## Slice
```
func append(slice []T, elements ...T) []T
```

```
a := []int{1, 2, 3}
b := []int{}
b = append(b, a...)
fmt.Println(b)
```