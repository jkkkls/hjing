package net

import (
	"fmt"
	"log"
	"os"

	"github.com/jkkkls/hjing/utils"
	kcp "github.com/xtaci/kcp-go"
)

func RunKCPServer(port int) error {
	utils.Info("启动kcp服务器", "port", port)
	if listener, err := kcp.ListenWithOptions(fmt.Sprintf(":%v", port), nil, 0, 0); err == nil {
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Println("AcceptKCP, err:", err)
				os.Exit(0)
			}
			s.SetNoDelay(1, 10, 2, 1)
			// s.SetMtu(512)
			// s.SetWindowSize(128, 128)
			s.SetStreamMode(true)
			s.SetWriteDelay(false)
			s.SetACKNoDelay(true)
			utils.Submit(func() {
				gater.handleConn(s, s.RemoteAddr().String(), "kcp")
			})
		}
	} else {
		return err
	}
}
