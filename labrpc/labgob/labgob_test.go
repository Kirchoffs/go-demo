package labgob

import (
    "bytes"
    "testing"
)

type AlphaStruct struct {
    IntKey    int
    IntVal    int
    StringKey string
    StringVal string
}

func TestGobFunctionality(t *testing.T) {
    Register(AlphaStruct{})

    num := 42
    str := "hello"
    alphaStruct := AlphaStruct{IntKey: num, IntVal: num, StringKey: str, StringVal: str}

    var buf bytes.Buffer

    encoder := NewEncoder(&buf)
    encoder.Encode(num)
    encoder.Encode(str)
    encoder.Encode(alphaStruct)

    decoder := NewDecoder(&buf)
    var numDec int
    var strDec string
    var alphaStructDec AlphaStruct
    decoder.Decode(&numDec)
    decoder.Decode(&strDec)
    decoder.Decode(&alphaStructDec)

    if num != numDec {
        t.Errorf("Expected %d, got %d", num, numDec)
    }

    if str != strDec {
        t.Errorf("Expected %s, got %s", str, strDec)
    }

    if alphaStruct != alphaStructDec {
        t.Errorf("Expected %v, got %v", alphaStruct, alphaStructDec)
    }
}

type BetaStruct struct {
    Yes bool
    no  bool
}

func TestLabGobCapitalWarning(t *testing.T) {
    initialErrorCount := errorCount

    Register(BetaStruct{})

    var buf bytes.Buffer

    encoder := NewEncoder(&buf)
    decoder := NewDecoder(&buf)

    encoder.Encode(BetaStruct{Yes: true, no: true})

    var betaStructDec BetaStruct
    decoder.Decode(&betaStructDec)

    if errorCount != initialErrorCount+1 {
        t.Errorf("Expected %d, got %d", initialErrorCount+1, errorCount)
    }

    if !betaStructDec.Yes {
        t.Errorf("Expected true, got false")
    }

    if betaStructDec.no {
        t.Errorf("Expected false, got true")
    }
}

type GammaStruct struct {
    IntKey int
}

func TestLabGobNonDefaultWarning(t *testing.T) {
    initialErrorCount := errorCount

    Register(GammaStruct{})

    var buf bytes.Buffer

    encoder := NewEncoder(&buf)
    decoder := NewDecoder(&buf)

    encoder.Encode(GammaStruct{IntKey: 42})

    var gammaStructDec GammaStruct
    gammaStructDec.IntKey = 89
    decoder.Decode(&gammaStructDec)

    if errorCount != initialErrorCount+1 {
        t.Errorf("Expected %d, got %d", initialErrorCount+1, errorCount)
    }
}
