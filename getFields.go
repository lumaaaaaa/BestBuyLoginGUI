package main

import (
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
)

var (
	alphaRegex       = regexp.MustCompile(`^[0-9]+_A_.+$`)
	passRegex        = regexp.MustCompile(`^\d+_X_.+$`)
	jsonPayloadRegex = regexp.MustCompile(`initData = \{.+\}`)
)

func reverseArray(input []string) []string {
	if len(input) == 0 {
		return input
	}
	return append(reverseArray(input[1:]), input[0])
}

func atob(input string) string {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return ""
	}
	return string(data)
}

func getPasswordCode(codeList []string) (string, error) {
	correct := ""
	for i, a := range codeList {
		a = atob(a)
		if passRegex.MatchString(a) {
			correct = codeList[i]
			break
		}
	}
	if correct == "" {
		return "", errors.New("unable to find correct codeList code")
	}
	return correct, nil
}

func getAlphaCode(alphaList []string) (string, error) {
	correct := ""
	for i, a := range alphaList {
		a = strings.Join(reverseArray(strings.Split(a, "")), "")
		a = atob(a)
		if alphaRegex.MatchString(a) {
			correct = alphaList[i]
			break
		}
	}
	if correct == "" {
		return "", errors.New("unable to find correct alpha")
	}
	return correct, nil
}
