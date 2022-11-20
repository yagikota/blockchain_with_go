package common

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"time"
)

func IsFoundHost(host string, port int) bool {
	target := net.JoinHostPort(host, strconv.Itoa(port))

	// https://christina04.hatenablog.com/entry/go-timeouts
	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		fmt.Printf("%s %v\n", target, err)
		return false
	}
	fmt.Printf("%s found\n", target)
	return true
}

// 疑問符は直前の表現が0個か1個あることを示す
// 0～9 または 10～99 または 100～199([01]?[0-9][0-9]?) または
// 200～249(2[0-4][0-9]) または
// 250～255 25[0-5]
var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

// IPアドレスの下１桁を走査(startIP~endIP)
// さらにポートも走査(startPort~endPort)
func FindNeighbors(myHost string, myPort int, startIp int, endIp int, startPort int, endPort int) []string {
	address := net.JoinHostPort(myHost, strconv.Itoa(myPort))

	// "127.0.0.1"
	// "127.0.0."
	// "0."
	// "1"
	// 切り出した文字列全体と、正規表現を「()」で括った単位ずつの個数分の要素です。
	// ()内に記載されていないものは、全体には含まれますが、単位要素には含まれません。
	// https://tech-up.hatenablog.com/entry/2018/12/04/224814
	m := PATTERN.FindStringSubmatch(myHost)
	if m == nil {
		return nil
	}
	prefixHost := m[1]
	lastIp, _ := strconv.Atoi(m[len(m)-1])
	neighbors := make([]string, 0, (endPort-startPort+1)*(endIp-startIp+1))

	for ip := startIp; ip <= endIp; ip += 1 {
		for port := startPort; port <= endPort; port += 1 {
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			if guessTarget != address && IsFoundHost(guessHost, port) {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}
	return neighbors
}
