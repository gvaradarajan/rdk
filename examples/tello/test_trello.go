package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/edaniels/golog"
	"go.viam.com/utils"
)

const reqAddr = "192.168.10.1:8889"

var logger = golog.NewDevelopmentLogger("tello")

func getResponse(cmdConn io.Reader) (string, error) {
	var buf [1518]byte
	_, err := cmdConn.Read(buf[0:])
	if err != nil {
		return "", err
	}
	return string(buf[:]), nil
}

func main() {
	utils.ContextualMain(mainWithArgs, logger)
}

func mainWithArgs(ctx context.Context, args []string, logger golog.Logger) error {
	exitCh := make(chan struct{}, 1)
	reqAddr, err := net.ResolveUDPAddr("udp", reqAddr)
	if err != nil {
		return err
	}
	respPort, err := net.ResolveUDPAddr("udp", ":9000")
	if err != nil {
		return err
	}
	cmdConn, err := net.DialUDP("udp", respPort, reqAddr)
	if err != nil {
		return err
	}
	defer func() {
		err := cmdConn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
	cmdLoop:
		for {
			select {
			case <-exitCh:
				logger.Info("closing response loop...")
				break cmdLoop
			default:
				// err := d.handleResponse(cmdConn)
				logger.Info("handling response...")
				resp, err := getResponse(cmdConn)
				if err != nil {
					fmt.Println("response parse error:", err)
				} else {
					logger.Info(resp)
				}
			}
		}
	}()
	fmt.Println("Enter command: ")

	for {
		var cmd string
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			cmd = scanner.Text()
		}
		if strings.Contains(cmd, "end") {
			logger.Info("closing...")
			exitCh <- struct{}{}
			break
		}
		logger.Info(fmt.Sprintf("Running command: %s", cmd))
		_, err := cmdConn.Write([]byte(cmd))
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
