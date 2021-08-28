package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

func TestHwTelnetClient_Connect_Failed(t *testing.T) {
	timeout := 1 * time.Millisecond
	client := NewTelnetClient("example.com", timeout, nil, nil)
	err := client.Connect()
	require.Error(t, err)
}

func TestHwTelnetClient_Send(t *testing.T) {
	client, l, in, _ := createClient(t)
	defer func() { require.NoError(t, l.Close()) }()

	require.NoError(t, client.Connect())
	defer func() { require.NoError(t, client.Close()) }()

	in.WriteString("hello\n")
	err := client.Send()
	require.NoError(t, err)

	require.Equal(t, "hello\n", readFromConn(t, l))
}

func TestHwTelnetClient_Receive(t *testing.T) {
	client, l, _, out := createClient(t)
	defer func() { require.NoError(t, l.Close()) }()

	require.NoError(t, client.Connect())
	defer func() { require.NoError(t, client.Close()) }()

	writeToConn(t, l, "world\n")

	err := client.Receive()
	require.NoError(t, err)
	require.Equal(t, "world\n", out.String())
}

func createClient(t *testing.T) (TelnetClient, net.Listener, *bytes.Buffer, *bytes.Buffer) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)

	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	timeout, err := time.ParseDuration("10s")
	require.NoError(t, err)

	client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)

	return client, l, in, out
}

func readFromConn(t *testing.T, l net.Listener) string {
	conn, err := l.Accept()
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer func() { require.NoError(t, conn.Close()) }()

	request := make([]byte, 1024)
	n, err := conn.Read(request)
	require.NoError(t, err)
	return string(request)[:n]
}

func writeToConn(t *testing.T, l net.Listener, s string) {
	conn, err := l.Accept()
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer func() { require.NoError(t, conn.Close()) }()

	n, err := conn.Write([]byte(s))
	require.NoError(t, err)
	require.NotEqual(t, 0, n)
}
