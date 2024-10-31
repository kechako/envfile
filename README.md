# envfile

Package envfile provides functionality to parse files containing environment variables in the format key=value.

## Install

```console
go get github.com/kechako/envfile@latest
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/kechako/envfile"
)

func main() {
	envs, err := envfile.ParseFile(".env")
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range envs.Envs() {
		fmt.Printf("%s = %v\n", key, value)
	}
}
```
