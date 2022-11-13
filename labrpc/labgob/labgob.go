package labgob

//
// For 6.584 labs.
// trying to send non-capitalized fields over RPC produces a range of
// misbehavior, including both mysterious incorrect computation and
// outright crashes. so this wrapper around Go's encoding/gob warns
// about non-capitalized field names.
//

import (
    "encoding/gob"
    "fmt"
    "io"
    "reflect"
    "sync"
    "unicode"
    "unicode/utf8"
)

var mu sync.Mutex
var errorCount int // for TestCapital
var checked map[reflect.Type]bool

type LabEncoder struct {
    gob *gob.Encoder
}

func NewEncoder(w io.Writer) *LabEncoder {
    enc := &LabEncoder{}
    enc.gob = gob.NewEncoder(w)
    return enc
}

func (enc *LabEncoder) Encode(e interface{}) error {
    checkFieldName(e)
    return enc.gob.Encode(e)
}

func (enc *LabEncoder) EncodeValue(value reflect.Value) error {
    checkFieldName(value.Interface())
    return enc.gob.EncodeValue(value)
}

type LabDecoder struct {
    gob *gob.Decoder
}

func NewDecoder(r io.Reader) *LabDecoder {
    dec := &LabDecoder{}
    dec.gob = gob.NewDecoder(r)
    return dec
}

func (dec *LabDecoder) Decode(e interface{}) error {
    checkFieldName(e)
    checkFieldValue(e)
    return dec.gob.Decode(e)
}

func Register(value interface{}) {
    checkFieldName(value)
    gob.Register(value)
}

func RegisterName(name string, value interface{}) {
    checkFieldName(value)
    gob.RegisterName(name, value)
}

func checkFieldName(value interface{}) {
    checkFieldNameInternal(reflect.TypeOf(value))
}

func checkFieldNameInternal(t reflect.Type) {
    mu.Lock()
    // only complain once, and avoid recursion.
    if checked == nil {
        checked = map[reflect.Type]bool{}
    }
    if checked[t] {
        mu.Unlock()
        return
    }
    checked[t] = true
    mu.Unlock()

    k := t.Kind()
    switch k {
    case reflect.Struct:
        for i := 0; i < t.NumField(); i++ {
            f := t.Field(i)
            rune, _ := utf8.DecodeRuneInString(f.Name)
            if !unicode.IsUpper(rune) {
                // ta da
                fmt.Printf("labgob error: lower-case field %v of %v in RPC or persist/snapshot will break your Raft\n",
                    f.Name, t.Name())
                mu.Lock()
                errorCount += 1
                mu.Unlock()
            }
            checkFieldNameInternal(f.Type)
        }
        return
    case reflect.Slice, reflect.Array, reflect.Ptr:
        checkFieldNameInternal(t.Elem())
        return
    case reflect.Map:
        checkFieldNameInternal(t.Elem())
        checkFieldNameInternal(t.Key())
        return
    default:
        return
    }
}

// warn if the value contains non-default values,
// as it would if one sent an RPC but the reply
// struct was already modified. if the RPC reply
// contains default values, GOB won't overwrite
// the non-default value with the default value.
func checkFieldValue(value interface{}) {
    if value == nil {
        return
    }
    checkFieldValueInternalWithDepth(reflect.ValueOf(value), 3, "")
}

func checkFieldValueInternalWithDepth(value reflect.Value, depth int, name string) {
    if depth == 0 {
        return
    }

    if !value.CanInterface() {
        return
    }

    t := value.Type()
    k := t.Kind()

    switch k {
    case reflect.Struct:
        for i := 0; i < t.NumField(); i++ {
            field := value.Field(i)
            fieldName := t.Field(i).Name
            if name != "" {
                fieldName = name + "." + fieldName
            }
            checkFieldValueInternalWithDepth(field, depth-1, fieldName)
        }
        return
    case reflect.Ptr:
        if value.IsNil() {
            return
        }
        checkFieldValueInternalWithDepth(value.Elem(), depth-1, name)
        return
    case reflect.Bool,
        reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
        reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
        reflect.Uintptr, reflect.Float32, reflect.Float64,
        reflect.String:
        if !reflect.DeepEqual(reflect.Zero(t).Interface(), value.Interface()) {
            mu.Lock()
            if errorCount < 1 {
                what := name
                if what == "" {
                    what = t.Name()
                }
                // this warning typically arises if code re-uses the same RPC reply
                // variable for multiple RPC calls, or if code restores persisted
                // state into variable that already have non-default values.
                fmt.Printf("labgob warning: Decoding into a non-default variable/field %v may not work\n", what)
            }
            errorCount += 1
            mu.Unlock()
        }
        return
    }
}
