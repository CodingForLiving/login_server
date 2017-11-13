package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
)

// 返回失败信息
func retFail(status int, info string, w http.ResponseWriter){
	m := map[string]interface{}{}
	m["status"] = status
	m["info"] = info
	retstr, _ := json.Marshal(m)
	fmt.Fprint(w,string(retstr))
	fmt.Print(string(retstr))
}

// 验证成功，存库
func OnAuthResult(w http.ResponseWriter, m map[string]interface{}) {
	if m == nil {
		retFail(0, "验证失败", w)
		return
	}
	fmt.Println(m)
	var acc = &Account{
		Id: 0,
		Openid: m["openid"].(string),
		Nickname: m["nickname"].(string),
		Sex: int(m["sex"].(float64)),
		Province: m["province"].(string),
		City: m["city"].(string),
		Country: m["country"].(string),
		Headimgurl: m["headimgurl"].(string),
		Unionid: m["unionid"].(string),
	}

	// 存库
	msg := db.OnWxLogin(acc)

	if msg != "success" {
		retFail(0, msg, w)
		return
	}

	ret := map[string]interface{}{}
	m["status"] = 1
	m["info"] = "登陆成功，返回登陆信息"
	info := map[string]interface{}{}
	info["unionid"] = acc.Unionid
	info["userid"] = acc.Id
	info["refresh_token"] = m["refreshtoken"].(string)
	info["dynamicpass"] = "fd"
	ret["info"] = info

	retstr, _ := json.Marshal(ret)
	fmt.Fprint(w,string(retstr))
	fmt.Print(string(retstr))
}

func handlerWxLogin(w http.ResponseWriter, r *http.Request){
	if len(r.Form["code"]) == 0 {
		return
	}
	var code = r.Form["code"][0]
	var result = wxAuthorize(code)
	OnAuthResult(w, result)
}

func handlerRefreshToken(w http.ResponseWriter, r *http.Request){
	log.Println("refreshtoken")
	if len(r.Form["refresh_token"]) != 1 {
		return
	}
	var token = r.Form["refresh_token"][0]
	fmt.Println(token)
	var result = wxRefreshTokenAuth(token)
	OnAuthResult(w, result)
}

func handlerGetGameSvr(w http.ResponseWriter, r *http.Request){
	m := map[string]interface{}{}
	m["status"] = 1
	m["info"] = "网关获取成功"
	api := map[string]interface{}{}
	api["gameserver"] = "192.168.1.20"
	api["uid"] = 111
	m["api"] = api
	retstr, _ := json.Marshal(m)
	fmt.Fprint(w,string(retstr))
	fmt.Print(string(retstr))
}

// 注册账号
func handlerRegister(w http.ResponseWriter, r *http.Request){
	if len(r.Form["account"]) != 1 {
		retFail(1, "注册账号信息不全：缺少用户名", w)
		return
	}
	if len(r.Form["password"]) != 1 {
		retFail(1, "注册账号信息不全：缺少密码", w)
		return
	}

	acc := &Account{
		Id: 0,
		Account: r.Form["account"][0],
		Password: r.Form["password"][0],
	}
	db.OnRegister(acc)
	m := map[string]interface{}{}
	m["id"] = acc.Id
	retstr, _ := json.Marshal(m)
	fmt.Fprint(w,string(retstr))
}

// 验证账号
func handlerLogin(w http.ResponseWriter, r *http.Request){
	if len(r.Form["account"]) != 1 {
		retFail(1, "验证账号信息不全：缺少用户名", w)
		return
	}
	if len(r.Form["password"]) != 1 {
		retFail(1, "验证账号信息不全：缺少密码", w)
		return
	}

	acc := &Account{
		Id: 0,
		Account: r.Form["account"][0],
		Password: r.Form["password"][0],
	}
	db.OnLogin(acc)

	m := map[string]interface{}{}
	m["id"] = acc.Id
	retstr, _ := json.Marshal(m)
	fmt.Fprint(w,string(retstr))
}

func checkSign(w http.ResponseWriter, r *http.Request) bool{
	fmt.Println(r.Header)
	fmt.Println(r.Header.Get("Sign"))
	cookie, err := r.Cookie("data")
	if err != nil {
		return false
	}
	fmt.Println("sign: ", cookie)
	return true
}

func handler(w http.ResponseWriter, r *http.Request){
	log.Println("handler request from ", r.RemoteAddr)
	r.ParseForm()
	checkSign(w, r)
	var param = r.Form[""][0]
}