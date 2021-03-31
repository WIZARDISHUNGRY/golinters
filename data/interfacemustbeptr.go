package testdata

//args: -Einterfacemustbeptr

import (
	"encoding/json"

	jason "encoding/json"
)

const data = `{"foo": "bar"}`

var dst = make(map[string]string)

type fancy struct {
	dst map[string]string
}

var myFancyDst fancy = fancy{dst: map[string]string{}}

// var fancyDst = {
// 	i int
// 	dst map[string]string
// }struct{
// 	dst: dst
// }

type myAlias = map[string]string
type myPtrAlias = *map[string]string

func Fail_UnmarshalMap() error {
	return json.Unmarshal([]byte(data), dst)
}

func Pass_MarshalMap_AddressOperator() error {
	return json.Unmarshal([]byte(data), &dst)
}

func Fail_UnmarshalMap_AliasedPackage() error {
	return jason.Unmarshal([]byte(data), dst)
}

func Fail_UnmarshalMap_Closure() error {
	return func() error {
		return json.Unmarshal([]byte(data), dst)
	}()
}

func Fail_UnmarshalMap_NamedClosure() error {
	f := func() error {
		return json.Unmarshal([]byte(data), dst)
	}
	return f()
}

func Fail_UnmarshalMap_Copy() error {
	myDst := dst
	return json.Unmarshal([]byte(data), myDst)
}

func Pass_UnmarshalMap_CopyAddressOperator() error {
	myDst := dst
	return json.Unmarshal([]byte(data), &myDst)
}

func Pass_UnmarshalMap_CreatePointer() error {
	myDst := &dst
	return json.Unmarshal([]byte(data), myDst)
}

func Pass_UnmarshalMap_TypeAlias() error {
	myDst := myAlias(dst)
	return json.Unmarshal([]byte(data), myDst)
}

func Fail_UnmarshalMap_TypeAlias() error {
	myDst := myAlias(dst)
	return json.Unmarshal([]byte(data), myDst)
}

func Pass_UnmarshalMap_AddressOfTypeAlias() error {
	myDst := myAlias(dst)
	return json.Unmarshal([]byte(data), &myDst)
}

func Pass_UnmarshalMap_PtrTypeAlias() error {
	myDst := myPtrAlias(&dst)
	return json.Unmarshal([]byte(data), myDst)
}

func addressOf(s map[string]string) *map[string]string {
	return &s
}

func Pass_UnmarshalMap_FunctionCall() error {
	return json.Unmarshal([]byte(data), addressOf(dst))
}

func identityOf(s map[string]string) map[string]string {
	return s
}

func Fail_UnmarshalMap_FunctionCall() error {
	return json.Unmarshal([]byte(data), identityOf(dst))
}

func Fail_UnmarshalMap_ConstExpr() error {
	return json.Unmarshal([]byte(data), map[string]string{})
}

func Fail_UnmarshalMap_Parens() error {
	return json.Unmarshal([]byte(data), (dst))
}

func Pass_UnmarshalMap_Parens() error {
	return json.Unmarshal([]byte(data), (&dst))
}

func Pass_UnmarshalMap_FancyDst() error {
	return json.Unmarshal([]byte(data), &myFancyDst.dst)
}

func Fail_UnmarshalMap_FancyDst() error {
	return json.Unmarshal([]byte(data), myFancyDst.dst)
}
