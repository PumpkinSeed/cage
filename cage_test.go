package cage

import (
	"fmt"
	"testing"
)

func TestStart(t *testing.T) {
	c := Start()

	fmt.Println("test")
	fmt.Println("lofasz")
	fmt.Println("test")

	Stop(c)

	fmt.Println(c.Data)
}