package testdata

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTestData(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Error(t, c.function())
		})
	}
}

// grep  'func [PF]' ./test/testdata/interfacemustbeptr.go | sed -e "s/.*func \(.*\)().*/{\1, \"\1\"},/"
var cases = []struct {
	function func() error
	name     string
	success  bool
}{
	{Fail_UnmarshalMap, "Fail_UnmarshalMap", false},
	{Pass_MarshalMap_AddressOperator, "Pass_MarshalMap_AddressOperator", true},
	{Fail_UnmarshalMap_AliasedPackage, "Fail_UnmarshalMap_AliasedPackage", false},
	{Fail_UnmarshalMap_Closure, "Fail_UnmarshalMap_Closure", false},
	{Fail_UnmarshalMap_NamedClosure, "Fail_UnmarshalMap_NamedClosure", false},
	{Fail_UnmarshalMap_Copy, "Fail_UnmarshalMap_Copy", false},
	{Pass_UnmarshalMap_CopyAddressOperator, "Pass_UnmarshalMap_CopyAddressOperator", true},
	{Pass_UnmarshalMap_CreatePointer, "Pass_UnmarshalMap_CreatePointer", true},
	{Pass_UnmarshalMap_TypeAlias, "Pass_UnmarshalMap_TypeAlias", true},
	{Fail_UnmarshalMap_TypeAlias, "Fail_UnmarshalMap_TypeAlias", false},
	{Pass_UnmarshalMap_AddressOfTypeAlias, "Pass_UnmarshalMap_AddressOfTypeAlias", true},
	{Pass_UnmarshalMap_PtrTypeAlias, "Pass_UnmarshalMap_PtrTypeAlias", true},
	{Pass_UnmarshalMap_FunctionCall, "Pass_UnmarshalMap_FunctionCall", true},
	{Fail_UnmarshalMap_FunctionCall, "Fail_UnmarshalMap_FunctionCall", false},
	{Fail_UnmarshalMap_ConstExpr, "Fail_UnmarshalMap_ConstExpr", false},
	{Fail_UnmarshalMap_Parens, "Fail_UnmarshalMap_Parens", false},
	{Pass_UnmarshalMap_Parens, "Pass_UnmarshalMap_Parens", true},
	{Pass_UnmarshalMap_FancyDst, "Pass_UnmarshalMap_FancyDst", true},
	{Fail_UnmarshalMap_FancyDst, "Fail_UnmarshalMap_FancyDst", false},
}
