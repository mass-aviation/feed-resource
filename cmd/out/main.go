package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintf(os.Stdout, `{"version":[],"metadata":{"name":"comment","value":"This resource has no output"}}`)
}
