package main

import (
	"container/list"
	"errors"
	"github.com/saisesai/doghole"
	"log"
	"net"
	"sync"
)

type Server struct {
	listener      net.Listener
	registry      list.List
	registryMutex sync.RWMutex
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() {
	var err error
	s.listener, err = net.Listen("tcp", *optionListenAddr)
	if err != nil {
		log.Panicln("failed to listen at", *optionListenAddr, ":", err)
	}
	log.Println("server listen at:", *optionListenAddr)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("failed to accept conn:", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println(conn.RemoteAddr(), "failed to close conn:", err)
		}
		log.Println(conn.RemoteAddr(), "closed")
	}()
	log.Println(conn.RemoteAddr(), "connected")

	var err error
	ctx := doghole.NewContext(conn)

	s.addContext(ctx)

	err = s.exchangeInfo(ctx)
	if err != nil {
		log.Println(ctx.Conn.RemoteAddr(), "failed to exchange info:", err)
		return
	}

	err = s.keepAlive(ctx)
	if err != nil {
		log.Println(ctx.Conn.RemoteAddr(), "failed to keep service alive:", err)
		return
	}
}

func (s *Server) addContext(ctx *doghole.Context) {
	s.registryMutex.Lock()
	s.registry.PushBack(ctx)
	s.registryMutex.Unlock()
}

func (s *Server) removeContext(ctx *doghole.Context) {
	s.registryMutex.Lock()
	for e := s.registry.Front(); e.Value.(*doghole.Context) != ctx; e = e.Next() {
		s.registry.Remove(e)
	}
	s.registryMutex.Unlock()
}

func (s *Server) exchangeInfo(ctx *doghole.Context) (err error) {
	// read info
	name, err := ctx.ReadLine()
	if err != nil {
		return err
	}
	ctx.Name = name
	// write info
	err = ctx.WriteLine(ctx.Conn.RemoteAddr().String())
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) keepAlive(ctx *doghole.Context) (err error) {
	for {
		// read ping
		var ping string
		ping, err = ctx.ReadLine()
		if err != nil {
			break
		}
		if ping != "ping" {
			err = errors.New("invalid ping: " + ping)
			break
		}
		// send pong
		err = ctx.WriteLine("pong")
		if err != nil {
			break
		}
	}
	return err
}
