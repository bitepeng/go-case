package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

// Query server const
const (
	IP_WHOIS_SERVER     = "whois.iana.org"
	DOMAIN_WHOIS_SERVER = "whois-servers.net"
	WHOIS_PORT          = "43"
	ARIN_SERVER         = "whois.arin.net"
)

// Whois do the whois query and returns whois info
func Whois(domain string, servers ...string) (result string, err error) {
	domain = strings.Trim(strings.TrimSpace(domain), ".")
	if domain == "" {
		err = fmt.Errorf("domain is empty")
		return
	}

	result, err = query(domain, servers...)
	if err != nil {
		return
	}

	token := "Registrar WHOIS Server:"
	if IsIpv4(domain) {
		token = "whois:"
	}

	start := strings.Index(result, token)
	if start == -1 {
		return
	}

	start += len(token)
	end := strings.Index(result[start:], "\n")
	server := strings.TrimSpace(result[start : start+end])
	if server == "" {
		return
	}

	tmpResult, err := query(domain, server)
	if err != nil {
		return
	}

	result += tmpResult

	return
}

// query do the query
func query(domain string, servers ...string) (result string, err error) {
	var server string
	if len(servers) == 0 || servers[0] == "" {
		if IsIpv4(domain) {
			server = IP_WHOIS_SERVER
		} else {
			domains := strings.Split(domain, ".")
			if len(domains) < 2 {
				err = fmt.Errorf("domain %s is invalid", domain)
				return
			}
			server = domains[len(domains)-1] + "." + DOMAIN_WHOIS_SERVER
		}
	} else {
		server = strings.ToLower(servers[0])
		if server == ARIN_SERVER {
			domain = "n + " + domain
		}
	}

	// 查询Whois域名
	hostUrl := net.JoinHostPort(server, WHOIS_PORT)
	conn, e := net.DialTimeout("tcp", hostUrl, time.Second*30)
	if e != nil {
		err = e
		return
	}
	defer func() { _ = conn.Close() }()
	_, _ = conn.Write([]byte(domain + "\r\n"))
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	// 查询
	buffer, e := ioutil.ReadAll(conn)
	if e != nil {
		err = e
		return
	}

	result = string(buffer)

	return
}

// IsIpv4 returns string is an ipv4 ip
func IsIpv4(ip string) bool {
	i := net.ParseIP(ip)
	return i.To4() != nil
}
