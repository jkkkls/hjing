package net

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jkkkls/hjing/utils"
)

func RunTCPServer(port int) error {
	utils.Info("启动tcp服务器", "port", port)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Accept tcp, err:", err)
			os.Exit(0)
		}
		utils.Submit(func() {
			gater.handleConn(conn, conn.RemoteAddr().String(), "tcp")
		})
	}
}
