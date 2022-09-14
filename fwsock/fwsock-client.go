package fwsock

import (
	"io"
	"os"
)

func (fws *FwSockClient) Close() {
	fws.Conn.Close()
}

func (fws *FwSockClient) SetChann(ch chan os.Signal) {
	fws.Ch = ch
}

func (fws *FwSockClient) Reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		println("Client got:", string(buf[0:n]))
	}
}
