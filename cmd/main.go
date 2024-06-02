package main

import (
	"strings"

	. "github.com/oneclickvirt/UnlockTests/defaultset"
	"github.com/oneclickvirt/UnlockTests/model"
)

func ShowResult(r model.Result) (s string) {
	switch r.Status {
	case model.StatusYes:
		s = Green("YES")
		if r.Info != "" {
			s += Green(" (" + r.Info + ")")
		}
		if r.Region != "" {
			s += Green(" (Region: " + strings.ToUpper(r.Region) + ")")
		}
		return s
	case model.StatusNetworkErr:
		return Red("NO") + Yellow(" (Network Err)")
	case model.StatusRestricted:
		s = Yellow("Restricted")
		if r.Info != "" {
			s += Yellow(" (" + r.Info + ")")
		}
		if r.Region != "" {
			s += Yellow(" (Region: " + strings.ToUpper(r.Region) + ")")
		}
		return s
	case model.StatusErr:
		s = Yellow("Error")
		if r.Err != nil {
			s += ": " + r.Err.Error()
		}
		return s
	case model.StatusNo:
		s = Red("NO")
		if r.Info != "" {
			s += Yellow(" (" + r.Info + ")")
		}
		if r.Region != "" {
			s += Yellow(" (Region: " + strings.ToUpper(r.Region) + ")")
		}
		return s
	case model.StatusBanned:
		s = Red("Banned")
		if r.Info != "" {
			s += Yellow(" (" + r.Info + ")")
		}
		return s
	case model.StatusUnexpected:
		s = Purple("Unexpected")
		if r.Err != nil {
			s += ": " + r.Err.Error()
		}
		return s
	default:
		return
	}
}
