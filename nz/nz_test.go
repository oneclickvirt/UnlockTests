package nz

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/utils"
	"testing"
)

func Test(t *testing.T) {
	req, _ := utils.ParseInterface("", "", "tcp4")

	res := NeonTV(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = SkyGO(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = ThreeNow(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = MaoriTV(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)
}
