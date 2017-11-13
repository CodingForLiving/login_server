package main

type Account struct {
	Id int				`json:"id"`
	Account string		`json:"account"`
	Password string		`json:"password"`
	Openid string		`json:"openid"`
	Nickname string		`json:"nickname"`
	Sex int				`json:"sex"`
	Province string		`json:"province"`
	City string			`json:"city"`
	Country string		`json:"country"`
	Headimgurl string	`json:"headimgurl"`
	Unionid string		`json:"unionid"`
}
