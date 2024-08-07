package utils

import (
	"fmt"
	"testing"
)

func TestDns(t *testing.T) {
	fmt.Println(get_nameserver_from_resolv())
}
