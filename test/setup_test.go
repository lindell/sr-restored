package test_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/lindell/sr-restored/run"
)

func setup(ctx context.Context, t *testing.T) *url.URL {
	port, err := getFreePort()
	if err != nil {
		t.Fatal(err)
	}
	addr := fmt.Sprintf("localhost:%d", port)

	go func() {
		if err := run.Run(ctx, run.Config{
			ServerAddr: addr,
		}); err != nil {
			t.Error(err)
		}
	}()

	u, err := url.Parse("http://" + addr)
	if err != nil {
		t.Fatal(err)
	}

	if err := waitForHTTP(u.String(), 5*time.Second); err != nil {
		t.Fatal(err)
	}

	return u
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func waitForHTTP(url string, maxWait time.Duration) error {
	before := time.Now()

	for {
		if time.Since(before) > maxWait {
			return errors.New("could not make successful http request")
		}

		_, err := http.Get(url)
		if err == nil {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}
}
