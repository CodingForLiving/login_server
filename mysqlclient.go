package main

import (
	"log"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
 )

type MysqlClient struct {
	db *sql.DB
	stmt *sql.Stmt
	addr string
	accounts []*Account
	c chan *Account
	closeChan chan bool
}

func (this *MysqlClient) init(addr string) {
	this.addr = addr
	this.accounts = []*Account{}
	this.c = make(chan *Account, 0)
	this.closeChan = make(chan bool, 0)
	this.connect()
}

// 连接
func (this *MysqlClient) connect(){
	db, err := sql.Open("mysql", this.addr)
    if err != nil {
        return
	}

	stmt, err1 := db.Prepare(`INSERT INTO tb_account(id, account_name) VALUES (?, ?)`)
	if err1 != nil {
		this.close()
		return
	}

	this.db = db
	this.stmt = stmt
}

func (this *MysqlClient)close(){
	if this.stmt != nil {
		this.stmt.Close()
		this.stmt = nil
	}

	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *MysqlClient)saveAccount(account *Account){
	ret, err := this.stmt.Exec(account.Id, account.Openid)
	if err != nil {
		log.Printf("insert data error: %v\n", err)
		return
	}
	if lastInsertId, err := ret.LastInsertId(); nil == err {
		log.Println("lastInsertId:", lastInsertId)
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err {
		log.Println("rowsAffected:", rowsAffected)
	}
}

// 批量存储
func (this *MysqlClient)save(){
	log.Println("批量处理mysql存储请求\n")
	var newAccounts = []*Account{}
	for _,acc := range this.accounts {
		_, err := this.stmt.Exec(acc.Id, acc.Openid)
		if err != nil {
			// 如果是断开连接

			// 数据错误
			newAccounts = append(newAccounts, acc)
			log.Printf("insert data error: %v\n", err)
		}
	}
	this.accounts = newAccounts
}

// 数据库协程
func (this *MysqlClient)loop(){
	timer := time.NewTimer(time.Second * 60)
	for {
        select {
        case acc := <- this.c:
            this.accounts = append(this.accounts, acc)
		case <- timer.C:
			this.save()
		case <- this.closeChan:
			this.Close()
            break
        }
    }
}

func (this *MysqlClient)NotifyClose(){
	this.closeChan <- true
}

func (this *MysqlClient)Close(){
	log.Println("mysql退出")
	this.save()
	this.db.Close()
}