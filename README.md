# Another Go validation library

# Usage

```go
package main

import (
	"fmt"

	"github.com/a-was/go-validator"
)

type MyStruct struct {
	Email string `flags:"required"`
	Name  string `regex:"^[a-z]+$"`
	Min   int    `min:"10" max:"20"`
	Max   int    `min:"10" max:"20"`
	Env   string `env:"MYAPP_SETTING" default:"default_value"`
}

func main() {
	myStruct := MyStruct{
		Email: "",
		Name:  "invalid-regex@example.com",
		Min:   5,
		Max:   30,
		Env:   "",
	}
	err := validator.Validate(&myStruct)
	fmt.Println(err)
	// Output:
	// Email: required value not filled
	// Name: invalid value: invalid-regex@example.com does not match regex ^[a-z]+$
	// Min: invalid value: 5, minimum value is 10
	// Max: invalid value: 30, maximum value is 20
}
```

# Avaliable tags

Tags can be combined with each other

## `min`
Takes integer as value <br />
For numeric types (int, uint, etc...) it checks if property value is **grater or equal** than value specyfied in tag <br />
For string, map, slice types it checks len

## `max`
Same as `min`, but checks if property value if **lower or equal**

## `regex`
Takes regex string as value <br />
Checks if string matches provided regex

## `flags`
Takes list of supported flags as value
### Supported flags
- `required` - means that property value cannot be zero-value or nil

## `env`
Takes environment variable name as value <br />
If property value is zero-value or nil, it sets its value to the value of the environment variable <br />

## `default`
Takes default value as value <br />
If property value is zero-value or nil, it sets its value to value specyfied in tag
