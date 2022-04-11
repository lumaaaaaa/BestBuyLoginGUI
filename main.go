package main

import (
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/valyala/fasthttp"
)

var (
	output = "Waiting for user...\n"
	pass   = ""
	email  = ""
)

func calllogin() {
	go login()
}

func callexit() {
	output += "-------------------\nShutting down...\n"
	go exit()
}

func login() {
	output += "-------------------\nInitializing login flow...\n"
	var req = fasthttp.AcquireRequest()
	var resp = fasthttp.AcquireResponse()
	output += "Getting token...\n"
	req.Header.SetMethod("GET")
	req.SetRequestURI("https://www.bestbuy.com/identity/global/signin?source=SWlsmVErR3RTh%2FWm%2Bn1NlzlbIS4FVv2%2Fobiuso7qC%2BJCL%2F7u3okDc9JUb9v%2BW53ANa%2FnNxphU7%2F1aTWHjM9e9w%3D%3D")
	req.Header.Add("User-Agent", "Bby-Android/21.10.10 APPSTORE Mozilla/5.0 (Linux; Android 7.1.2; G011A Build/N2G48H; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/68.0.3440.70 Mobile Safari/537.36")
	err := fasthttp.Do(req, resp)
	if err != nil {
		return
	}
	var token = string(resp.Header.Peek("Location"))
	token = token[46:]
	output += "Obtained token, " + token + "\n"
	req.SetRequestURI("https://www.bestbuy.com/identity/signin?token=" + token)
	var mailfield string
	var alphas []string
	var passfield []string
	err = fasthttp.Do(req, resp)
	if err != nil {
		return
	}
	var respbody = string(resp.Body())
	mailfield = strings.Split(strings.Split(respbody, "emailFieldName\":\"")[1], "\",")[0]
	var alphahold = strings.Split(strings.Split(respbody, "alpha\":[")[1], "]")[0]
	var passhold = strings.Split(strings.Split(respbody, "codeList\":[")[1], "]")[0]
	alphahold = strings.ReplaceAll(alphahold, "\"", "")
	passhold = strings.ReplaceAll(passhold, "\"", "")
	alphas = strings.Split(alphahold, ",")
	passfield = strings.Split(passhold, ",")
	req.SetRequestURI("https://www.bestbuy.com/identity/authenticate")
	req.Header.SetMethod("POST")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:97.0) Gecko/20100101 Firefox/97.0")
	req.Header.Set("Host", "www.bestbuy.com")
	req.Header.Set("Origin", "https://www.bestbuy.com")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Seatac", "accept")
	req.Header.Set("X-Touch-Id", "true")
	req.Header.Set("Referer", "https://www.bestbuy.com/identity/signin?token="+token)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("X-Requested-With", "com.bestbuy.android")
	req.Header.Set("Cookie", "ZPLANK=")
	alpha, err := getAlphaCode(alphas)
	if err != nil {
		output += "Unable to get correct alpha. This could happen if you're restricted from using bestbuy on this IP address."
		return
	}
	passField, err := getPasswordCode(passfield)
	if err != nil {
		output += "Unable to get correct password field. This could happen if you're restricted from using bestbuy on this IP address."
		return
	}
	req.SetBody([]byte("{\"token\":\"" + token + "\",\"loginMethod\":\"UID_PASSWORD\",\"flowOptions\":\"0000000000000000\",\"enrollBiometric\":true,\"alpha\":\"" + alpha + "\",\"Salmon\":\"FA7F2\",\"" + passField + "\":\"" + pass + "\",\"" + mailfield + "\":\"" + email + "\"}"))
	err = fasthttp.Do(req, resp)
	if err != nil {
		return
	}
	b := string(resp.Body())
	if strings.Contains(b, "success") {
		output += "Logged in successfully!\n"
	} else if strings.Contains(b, "expired") {
		output += "Not this one...\n"
	} else if strings.Contains(b, "failed") {
		output += "An error occurred checking this combo.\n"
	} else if strings.Contains(b, "failure") {
		output += "The provided credentials are invalid.\n"
	} else if strings.Contains(b, "stepUpRequired") {
		output += "2FA required. Reset in your browser and retry.\n"
	} else {
		output += "Unknown output: " + b
	}
}

func exit() {
	exit()
}

func loop() {
	g.SingleWindow().Layout(
		g.Label("BestBuy Login"),
		g.Separator(),
		g.InputTextMultiline(&output).Size(-1, 220),
		g.Row(
			g.InputText(&email).Hint("mail").Size(188),
			g.InputText(&pass).Hint("pass").Size(188),
		),
		g.Row(
			g.Button("Login").OnClick(calllogin),
			g.Button("Exit").OnClick(callexit),
		),
	)
}

func main() {
	g.SetDefaultFont("Consola", 12)
	wnd := g.NewMasterWindow("BestBuy ---", 400, 300, g.MasterWindowFlagsNotResizable)
	wnd.Run(loop)
}
