package tool

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func InitLog() {
	logFileName := flag.String("log", "./disk.log", "Log file name")
	flag.Parse()

	// 设置存储的路径
	logFile, logErr := os.OpenFile(*logFileName, os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "disk start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
