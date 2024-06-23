package au

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/utils"
	"testing"
)

func Test(t *testing.T) {
	req, _ := utils.ParseInterface("", "", "tcp4")

	res := Au7plus(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = ABCiView(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = Binge(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = Channel10(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = Channel9(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = OptusSports(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = SBSonDemand(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = Stan(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)

	res = KayoSports(req)
	if res.Err != nil {
		fmt.Println(res.Err)
	}
	fmt.Println(res.Name, ": ", res.Status, res.Region, res.UnlockType)
}
