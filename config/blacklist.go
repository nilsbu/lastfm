package config

import "strings"

func Blacklist() map[string]interface{} {
	list := []string{
		"20s", "30s", "40s", "50s", "60s", "70s", "80s", "90s", "00s", "10s",
		"beautiful", "sexy", "love",
		"comedy",
		"cover",
		"female vocalist", "female vocalists", "male vocalist", "male vocalists",
		"hip hop", "rap",
		"oldies",
		"psychedelic",
		"seen live",
		"soundtrack",
	}

	blacklist := map[string]interface{}{}
	for _, item := range list {
		blacklist[item] = nil
	}

	for k, v := range Countries {
		blacklist[k] = nil
		blacklist[strings.ToLower(v)] = nil
	}

	return blacklist
}
