package doghole

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type Context struct {
	Name string
	Conn net.Conn
	WR   *bufio.ReadWriter
}

func NewContext(conn net.Conn) *Context {
	return &Context{
		Conn: conn,
		WR:   bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
	}
}

func (ctx *Context) ReadLine() (str string, err error) {
	str, err = ctx.WR.ReadString('\n')
	if err != nil {
		return "", err
	}
	if *optionDebug {
		log.Println(ctx.Conn.RemoteAddr(), "->", []byte(str))
	}
	str = strings.ReplaceAll(str, "\r\n", "")
	str = strings.ReplaceAll(str, "\n", "")
	//str = strings.ReplaceAll()
	return str, err
}

func (ctx *Context) WriteLine(str string) (err error) {
	_, err = ctx.WR.WriteString(str + "\n")
	if err != nil {
		return err
	}
	err = ctx.WR.Flush()
	if err != nil {
		return err
	}
	if *optionDebug {
		log.Println(ctx.Conn.RemoteAddr(), "<-", []byte(str+"\n"))
	}
	return nil
}
