package main

import (
	"fmt"
)

func main() {
	size := 32 << (^uint(0) >> 63)
	fmt.Println("Size of uint on this machine:", size, "bits.")
}
