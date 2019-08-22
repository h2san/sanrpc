package client

import (
	"bufio"
	"crypto/tls"
	"net"
	"time"

	log "github.com/hillguo/sanlog"
)

// ReaderBuffSize is used for bufio reader.
const ReaderBuffSize = 16 * 1024

type makeConnFn func(c *Client, network, address string) (net.Conn, error)

var makeConnMap = make(map[string]makeConnFn)

// Connect connects the server via specified network.
func (client *Client) Connect(network, address string) error {
	var conn net.Conn
	var err error

	switch network {
	case "unix":
		conn, err = newDirectConn(client, network, address)
	default:
		fn := makeConnMap[network]
		if fn != nil {
			conn, err = fn(client, network, address)
		} else {
			conn, err = newDirectConn(client, network, address)
		}
	}

	if err == nil && conn != nil {
		if client.option.ReadTimeout != 0 {
			_ = conn.SetReadDeadline(time.Now().Add(client.option.ReadTimeout))
		}
		if client.option.WriteTimeout != 0 {
			_ = conn.SetWriteDeadline(time.Now().Add(client.option.WriteTimeout))
		}

		client.Conn = conn
		client.r = bufio.NewReaderSize(conn, ReaderBuffSize)

		// start reading and writing since connected
		go client.input()

		if client.option.Heartbeat && client.option.HeartbeatInterval > 0 {
			//go client.heartbeat()
		}

	}

	return err
}

func newDirectConn(c *Client, network, address string) (net.Conn, error) {
	var conn net.Conn
	var tlsConn *tls.Conn
	var err error

	if c != nil && c.option.TLSConfig != nil {
		dialer := &net.Dialer{
			Timeout: c.option.ConnectTimeout,
		}
		tlsConn, err = tls.DialWithDialer(dialer, network, address, c.option.TLSConfig)
		//or conn:= tls.Client(netConn, &config)
		conn = net.Conn(tlsConn)
	} else {
		conn, err = net.DialTimeout(network, address, c.option.ConnectTimeout)
	}

	if err != nil {
		log.Warnf("failed to dial server: %v", err)
		return nil, err
	}

	if tc, ok := conn.(*net.TCPConn); ok {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(3 * time.Minute)
	}

	return conn, nil
}


