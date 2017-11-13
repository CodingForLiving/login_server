package main

import (
	"httplib"
)

const wechat_getaccesstoken_url  = "https://api.weixin.qq.com/sns/oauth2/access_token"
const wechat_getuserinfo_url = "https://api.weixin.qq.com/sns/userinfo"
const wechat_refreshtoken_url = "https://api.weixin.qq.com/sns/oauth2/refresh_token"

var(
	wxappkey string = ""
	wxappsecret string = ""
)

func wxGetAccesstoken(code string) map[string]interface{}{
	request:= httplib.Get(wechat_getaccesstoken_url)
	request.Param("appid",wxappkey)
	request.Param("secret",wxappsecret)
	request.Param("code",code)
	request.Param("grant_type","authorization_code")
	var response map[string]interface{}
	err := request.ToJson(&response)
	if err != nil {
		return nil
	}
	return response
}

func wxRefreshToken(token string)map[string]interface{}{
	request:= httplib.Get(wechat_refreshtoken_url)
	request.Param("appid",wxappkey)
	request.Param("grant_type","refresh_token")
	request.Param("refresh_token",token)
	var response map[string]interface{}
	err := request.ToJson(&response)
	if err != nil {
		return nil
	}
	return response
}

func wxGetUserinfo(accesstoken string, openid string) map[string]interface{}{
	request:= httplib.Get(wechat_getuserinfo_url)
	request.Param("access_token", accesstoken)
	request.Param("openid", openid)
	var response map[string]interface{}
	err := request.ToJson(&response)
	if err != nil {
		return nil
	}
	return response
}

func wxAuthorize(code string) map[string]interface{}{
	accesstokenResponse := wxGetAccesstoken(code)
	if accesstokenResponse == nil{
		return nil
	}
	_, ok := accesstokenResponse["errcode"]         //获取accesstoken接口返回错误码
	if ok {
		return nil
	}
	openid := accesstokenResponse["openid"].(string)
	accesstoken := accesstokenResponse["access_token"].(string)
	refresh_token := accesstokenResponse["refresh_token"].(string)
	getuserinfoResult := wxGetUserinfo(accesstoken, openid)
	if getuserinfoResult == nil {
		return nil
	}
	_, ok = getuserinfoResult["errcode"]           //获取用户信息接口返回错误码
	if ok {
		return nil
	}

	getuserinfoResult["refreshtoken"] = refresh_token
	return getuserinfoResult
}

func wxRefreshTokenAuth(token string) map[string]interface{}{
	refreshResponse := wxRefreshToken(token)
	if refreshResponse == nil{
		return nil
	}
	_, ok := refreshResponse["errcode"]         //获取accesstoken接口返回错误码
	if ok {
		return nil
	}
	openid := refreshResponse["openid"].(string)
	accesstoken := refreshResponse["access_token"].(string)
	refresh_token := refreshResponse["refresh_token"].(string)
	getuserinfoResult := wxGetUserinfo(accesstoken, openid)

	if getuserinfoResult == nil {
		return nil
	}
	_, ok = getuserinfoResult["errcode"]           //获取用户信息接口返回错误码
	if ok {
		return nil
	}

	getuserinfoResult["refreshtoken"] = refresh_token
	return getuserinfoResult
}
