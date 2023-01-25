package validator

import (
	"os"
	"strings"
	"testing"
)

func TestMin(t *testing.T) {
	type test struct {
		Min    int  `min:"3"`
		MinPtr *int `min:"3"`
	}
	x := &test{
		Min:    3,
		MinPtr: ptr(3),
	}
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	x.Min = 2
	x.MinPtr = ptr(2)
	if err := Validate(x); err == nil {
		t.Fail()
	}
	x.Min = 3
	x.MinPtr = nil
	if err := Validate(x); err != nil {
		t.Error(err)
	}
}

func TestMax(t *testing.T) {
	type test struct {
		Max    int  `max:"5"`
		MaxPtr *int `max:"5"`
	}
	x := &test{
		Max:    3,
		MaxPtr: ptr(3),
	}
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	x.Max = 6
	x.MaxPtr = ptr(6)
	if err := Validate(x); err == nil {
		t.Fail()
	}
	x.Max = 5
	x.MaxPtr = nil
	if err := Validate(x); err != nil {
		t.Error(err)
	}
}

func TestRegex(t *testing.T) {
	type test struct {
		Regex    string  `regex:"[a-z]+"`
		RegexPtr *string `regex:"[a-z]+"`
	}
	x := &test{
		Regex:    "aaaaa",
		RegexPtr: ptr("aaaaa"),
	}
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	x.Regex = "aaAaa"
	x.RegexPtr = ptr("aaAaa")
	if err := Validate(x); err == nil {
		t.Fail()
	}
	x.Regex = "aaaaa"
	x.RegexPtr = nil
	if err := Validate(x); err != nil {
		t.Error(err)
	}
}

func TestFlags(t *testing.T) {
	type test struct {
		Required    string  `flags:"required"`
		RequiredPtr *string `flags:"required"`
	}
	x := &test{
		Required:    "v",
		RequiredPtr: ptr("v"),
	}
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	x.Required = ""
	x.RequiredPtr = nil
	if err := Validate(x); err == nil {
		t.Error("RequiredPtr nil")
	}
	x.Required = "a"
	x.RequiredPtr = ptr("")
	if err := Validate(x); err == nil {
		t.Error("RequiredPtr zero value")
	}
}

func TestEnv(t *testing.T) {
	type test struct {
		FromEnv        string  `env:"FROM_ENV"`
		FromEnvPtr     *string `env:"FROM_ENV_PTR"`
		FromEnvBool    bool    `env:"FROM_ENV_BOOL"`
		FromEnvBoolPtr *bool   `env:"FROM_ENV_BOOL_PTR"`
	}
	x := &test{}
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	// string
	os.Setenv("FROM_ENV", "aaa")
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	if x.FromEnv != "aaa" {
		t.Error("invalid x.FromEnv value:", x.FromEnv)
	}
	os.Setenv("FROM_ENV_PTR", "aaa")
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	if x.FromEnvPtr == nil {
		t.Error("invalid x.FromEnvPtr value: nil")
	}
	if *x.FromEnvPtr != "aaa" {
		t.Error("invalid x.FromEnvPtr value:", x.FromEnvPtr)
	}
	// bool
	os.Setenv("FROM_ENV_BOOL", "true")
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	if !x.FromEnvBool {
		t.Error("invalid x.FromEnvBool value:", x.FromEnvBool)
	}
	os.Setenv("FROM_ENV_BOOL_PTR", "true")
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	if x.FromEnvBoolPtr == nil {
		t.Error("invalid x.FromEnvBool value: nil")
	}
	if !*x.FromEnvBoolPtr {
		t.Error("invalid x.FromEnvBoolPtr value:", x.FromEnvBoolPtr)
	}
}

func TestEnvOverride(t *testing.T) {
	type test struct {
		Value string `env:"FROM_ENV"`
	}
	x := &test{Value: "value"}
	os.Setenv("FROM_ENV", "aaa")
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	if x.Value != "aaa" {
		t.Error("Value not changed:", x.Value)
	}
}

func TestRequiredFromEnv(t *testing.T) {
	type test struct {
		Value string `flags:"required" env:"FROM_ENV"`
	}
	x := &test{}
	os.Setenv("FROM_ENV", "aaa")
	if err := Validate(x); err != nil {
		t.Error(err)
	}
}

func TestDefault(t *testing.T) {
	type test struct {
		String string `flags:"required" default:"123"`
		Int    int    `flags:"required" default:"123"`
		Bool   bool   `flags:"required" default:"true"`
	}
	x := &test{}
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	if x.String != "123" {
		t.Error("Invalid String default value:", x.String)
	}
	if x.Int != 123 {
		t.Error("Invalid Int default value:", x.Int)
	}
	if !x.Bool {
		t.Error("Invalid Bool default value:", x.Bool)
	}
}

func TestInner(t *testing.T) {
	type inner struct {
		MinMax   int    `min:"3" max:"4"`
		Required string `flags:"required"`
	}
	type test struct {
		Inner inner `flags:"required"`
	}
	x := &test{
		Inner: inner{
			MinMax:   3,
			Required: "a",
		},
	}
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	x.Inner.MinMax = 2
	x.Inner.Required = ""
	if err := Validate(x); err == nil {
		t.Fail()
	}
}

func TestInnerPtr(t *testing.T) {
	type inner struct {
		MinMax   int    `min:"3" max:"4"`
		Required string `flags:"required"`
	}
	type test struct {
		Inner *inner `flags:"required"`
	}
	x := &test{
		Inner: &inner{
			MinMax:   3,
			Required: "a",
		},
	}
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	x.Inner = nil
	if err := Validate(x); err == nil {
		t.Fail()
	}
}

func TestDefaultTags(t *testing.T) {
	type test struct {
		_ struct{} `flags:"required"`
		A string
		B string
	}
	x := &test{
		A: "a",
		B: "b",
	}
	if err := Validate(x); err != nil {
		t.Error(err)
	}
	x.A = ""
	x.B = ""
	if err := Validate(x); err == nil {
		t.Fail()
	}
}

func TestErrors(t *testing.T) {
	type test struct {
		Min   int    `min:"3"`
		Max   int    `max:"5"`
		Regex string `regex:"[a-z]+"`
	}
	x := &test{
		Max: 7,
	}

	err := Validate(x)
	if err == nil {
		t.Fail()
	}

	msg := err.Error()
	if !strings.Contains(msg, "\n") {
		t.Error("invalid multiline error message", msg)
	}
	if !strings.Contains(msg, "Min: invalid value") && !strings.Contains(msg, "minimum value is") {
		t.Error("invalid min error message", msg)
	}
	if !strings.Contains(msg, "Max: invalid value") && !strings.Contains(msg, "maximum value is") {
		t.Error("invalid max error message", msg)
	}
	if !strings.Contains(msg, "Regex: invalid value") && !strings.Contains(msg, "empty string") {
		t.Error("invalid empty regex error message", msg)
	}
	x.Regex = "A"
	err = Validate(x)
	msg = err.Error()
	if !strings.Contains(msg, "Regex: invalid value") && !strings.Contains(msg, "does not match regex") {
		t.Error("invalid regex match error message", msg)
	}
}
