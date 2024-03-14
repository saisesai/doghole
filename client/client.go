package main

import (
	"errors"
	"github.com/libp2p/go-reuseport"
	"github.com/saisesai/doghole"
	"io"
	"log"
	"net"
	"time"
)

type Client struct {
	ctx        *doghole.Context
	listener   net.Listener
	publicAddr string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Run() {
	var err error
	err = c.localListen()
	if err != nil {
		log.Panicln("failed to listen:", err)
	}
	go c.serviceServe()
	for {
		err = c.serverConnect()
		if err != nil {
			log.Panicln("failed to connect server:", err)
		}
		err = c.keepAlive()
		if err != nil {
			log.Panicln("failed to keep alive:", err)
		}
		time.Sleep(time.Second * 3)
	}
}

func (c *Client) localListen() (err error) {
	c.listener, err = reuseport.Listen("tcp", *optionLocalAddr)
	return err
}

func (c *Client) serverConnect() (err error) {
	conn, err := reuseport.Dial("tcp", *optionLocalAddr, *optionServerAddr)
	if err != nil {
		return err
	}
	c.ctx = doghole.NewContext(conn)
	c.ctx.Name = *optionName
	// send info
	err = c.ctx.WriteLine(*optionName)
	if err != nil {
		return err
	}
	// get public addr
	c.publicAddr, err = c.ctx.ReadLine()
	if err != nil {
		return err
	}
	log.Println("client public addr:", c.publicAddr)
	return nil
}

func (c *Client) keepAlive() (err error) {
	for {
		// send ping
		err = c.ctx.WriteLine("ping")
		if err != nil {
			break
		}
		// read pong
		var pong string
		pong, err = c.ctx.ReadLine()
		if err != nil {
			break
		}
		if pong != "pong" {
			err = errors.New("invalid pong: " + pong)
			break
		}
		time.Sleep(time.Second * 3)
	}
	return err
}

func (c *Client) serviceServe() {
	for {
		conn, err := c.listener.Accept()
		if err != nil {
			log.Println("failed to accept:", err)
			continue
		}
		go c.handleConn(conn)
	}
}

func (c *Client) handleConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println(conn.RemoteAddr(), "failed to close conn:", err)
		}
		log.Println(conn.RemoteAddr(), "closed")
	}()
	log.Println(conn.RemoteAddr(), "connected")
	// dial to local service
	localConn, err := net.Dial("tcp", *optionServiceAddr)
	if err != nil {
		log.Println("failed to dial to local server:", err)
		return
	}
	// copy
	go io.Copy(conn, localConn)
	io.Copy(localConn, conn)
}
