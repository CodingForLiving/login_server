package main

import (
	"os"
	"net/http"
	"log"
	"flag"
	"os/signal"
	"syscall"
)

var db = Db{}

func initLog(){
	logF := flag.String("log", "login.log", "Log file name")
	flag.Parse() //解析参数付给logF
	outfile, err := os.OpenFile(*logF, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666) //打开文件，若果文件不存在就创
	if err != nil {
		panic("init log failed!")
	}
	log.SetOutput(outfile)  //设置log的输出文件，不设置log输出默认为stdout
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) //设置答应日志每一行前的标志信息，这里设置了日期，打印时间，
}

func closeLog(){

}

func main(){
	initLog()

	if len(os.Args) != 2 {
		log.Fatalf("please identify a config file ！")
	}

	log.Println("begin to start server...")

	LoadConfig(os.Args[1])

	db.Init()

	quitChan := make(chan os.Signal)
    signal.Notify(quitChan,
        syscall.SIGINT,
        syscall.SIGTERM,
        syscall.SIGHUP,
	)

	http.HandleFunc("/", handler)
	server := &http.Server{
		Addr: config.HttpAddr,
		Handler: http.DefaultServeMux,
	}

	go func() {
        <-quitChan
		server.Close()
		db.Close()
    }()

	err := server.ListenAndServe()
	if err != nil {
		log.Println("http server error: ", err.Error())
		return
	}
}