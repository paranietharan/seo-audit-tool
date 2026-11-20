package scraper

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/netip"
	"time"
)

// Pin the connection to the address that passed validation. Chromium uses the
// same gate for navigations, redirects and page subresources.
func dialPublic(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip", host)
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		if !publicAddress(ip) {
			return nil, fmt.Errorf("destination is outside the public network")
		}
	}
	dialer := net.Dialer{Timeout: 10 * time.Second}
	for _, ip := range ips {
		conn, err := dialer.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
		if err == nil {
			return conn, nil
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}
	return nil, fmt.Errorf("could not connect to public destination")
}

var nonPublicRanges = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"), netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("198.18.0.0/15"), netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("64:ff9b::/96"), netip.MustParsePrefix("64:ff9b:1::/48"),
	netip.MustParsePrefix("2002::/16"), netip.MustParsePrefix("2001::/32"),
}

func publicAddress(ip netip.Addr) bool {
	ip = ip.Unmap()
	if !ip.IsGlobalUnicast() || ip.IsPrivate() {
		return false
	}
	for _, prefix := range nonPublicRanges {
		if prefix.Contains(ip) {
			return false
		}
	}
	return true
}

func startBrowserProxy(transport *http.Transport) (string, func(), error) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return "", nil, err
	}
	forward := httputil.ReverseProxy{Transport: transport, Rewrite: func(req *httputil.ProxyRequest) {
		req.Out.Host = req.Out.URL.Host
		req.Out.Header.Del("Proxy-Authorization")
	}}
	server := &http.Server{ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 30 * time.Second}
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodConnect {
			if r.URL.Scheme != "http" || r.URL.Host == "" {
				http.Error(w, "HTTP target required", 400)
				return
			}
			forward.ServeHTTP(w, r)
			return
		}
		upstream, err := dialPublic(r.Context(), "tcp", r.Host)
		if err != nil {
			http.Error(w, "Destination unavailable", 502)
			return
		}
		defer upstream.Close()
		client, buffer, err := http.NewResponseController(w).Hijack()
		if err != nil {
			return
		}
		defer client.Close()
		deadline := time.Now().Add(time.Minute)
		client.SetDeadline(deadline)
		upstream.SetDeadline(deadline)
		if _, err = buffer.WriteString("HTTP/1.1 200 Connection Established\r\n\r\n"); err != nil {
			return
		}
		if err = buffer.Flush(); err != nil {
			return
		}
		go func() { io.Copy(upstream, buffer.Reader); upstream.Close() }()
		io.Copy(client, upstream)
	})
	go server.Serve(listener)
	closeProxy := func() { server.Close(); transport.CloseIdleConnections() }
	return "http://" + listener.Addr().String(), closeProxy, nil
}
