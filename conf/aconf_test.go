package conf

import (
	"testing"
	"fmt"
)

func TestLoad(t *testing.T) {

	_, e := LoadFile("../assembly.json")

	t.Log(e)
}

func TestGetCurrentPath(t *testing.T) {

	fmt.Println()

}
