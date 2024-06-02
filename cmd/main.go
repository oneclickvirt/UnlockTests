package main

import (
	"os"
	"strings"
	"sync"
	"time"

	. "github.com/oneclickvirt/UnlockTests/defaultset"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	pb "github.com/schollz/progressbar/v3"
)

var (
	tot int64
	bar *pb.ProgressBar
	wg *sync.WaitGroup
	IPV4 = true
	IPV6 = true
)

func NewBar(count int64) *pb.ProgressBar {
	return pb.NewOptions64(
		count,
		pb.OptionSetDescription("testing"),
		pb.OptionSetWriter(os.Stderr),
		pb.OptionSetWidth(20),
		pb.OptionThrottle(100*time.Millisecond),
		pb.OptionShowCount(),
		pb.OptionClearOnFinish(),
		pb.OptionEnableColorCodes(true),
		pb.OptionSpinnerType(14),
	)
}

func ShowResult(r model.Result) (s string) {
	formatResult := func(colorFunc func(string) string, status string, r model.Result) string {
		s := colorFunc(status)
		if r.Info != "" {
			s += colorFunc(" (" + r.Info + ")")
		}
		if r.Region != "" {
			s += colorFunc(" (Region: " + strings.ToUpper(r.Region) + ")")
		}
		return s
	}
	switch r.Status {
	case model.StatusYes:
		return formatResult(Green, "YES", r)
	case model.StatusNetworkErr:
		return Red("NO") + Yellow(" (Network Err)")
	case model.StatusRestricted:
		return formatResult(Yellow, "Restricted", r)
	case model.StatusErr:
		s = Yellow("Error")
		if r.Err != nil {
			s += ": " + r.Err.Error()
		}
		return s
	case model.StatusNo:
		return formatResult(Red, "NO", r)
	case model.StatusBanned:
		s = Red("Banned")
		if r.Info != "" {
			s += Yellow(" (" + r.Info + ")")
		}
		return s
	case model.StatusUnexpected:
		s = Purple("Unknown")
		if r.Err != nil {
			s += ": " + r.Err.Error()
		}
		return s
	default:
		return ""
	}
}

func excute(F func(request *gorequest.SuperAgent) model.Result, request *gorequest.SuperAgent) {
	wg.Add(1)
	tot++
	go func() {
		defer wg.Done()
		res := F(request)
		bar.Describe(res.Name + " " + ShowResult(res))
		bar.Add(1)
	}()
}
