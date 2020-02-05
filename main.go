package main

// dStor Echo Server
// MIT License // Stephanie Sunshine // 01/20

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// LookupUUID Function queries a remote ip address and port for a uuid
func LookupUUID(ip net.IP, port uint64) error {
	address := fmt.Sprintf("[%s]:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)

	if err != nil {
		return err
	}

	buffer := make([]byte, 512)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

	bufio.NewReader(conn).Read(buffer)

	re := regexp.MustCompile("[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}")
	match := re.FindStringSubmatch(string(buffer))
	defer conn.Close()
	if len(match) != 1 {
		return nil
	}
	return fmt.Errorf("%v", match[0])
}

// Reverse Reverse a string primitive
func Reverse(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, rune := range s {
		n--
		runes[n] = rune
	}
	return string(runes[n:])
}

// Entry
func main() {

	var PORT string

	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	} else {
		PORT = "80"
	}

	http.HandleFunc("/address", AddressLookup)
	http.HandleFunc("/uuid/", UUIDLookup)

	fmt.Printf("Starting server at :%s...\n", PORT)

	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

// AddressLookup Client ip address responder
func AddressLookup(w http.ResponseWriter, r *http.Request) {
	res := strings.SplitN(Reverse(r.RemoteAddr), ":", -1)
	remoteIP := Reverse(strings.Join(res[1:], ":"))

	re := regexp.MustCompile("\\[|\\]")
	remoteIP = re.ReplaceAllString(remoteIP, "")

	fmt.Fprintf(w, "%v\n", remoteIP)
}

// UUIDLookup Client uuid responder
func UUIDLookup(w http.ResponseWriter, r *http.Request) {
	res := strings.SplitN(Reverse(r.RemoteAddr), ":", -1)
	remoteIP := Reverse(strings.Join(res[1:], ":"))

	re := regexp.MustCompile("\\[|\\]")
	remoteIP = re.ReplaceAllString(remoteIP, "")

	ip := net.ParseIP(remoteIP)

	path := strings.SplitN(string(r.URL.Path[1:]), "/", -1)
	port, _ := strconv.ParseUint(path[1], 10, 64)
	fmt.Fprintf(w, "%v\n", LookupUUID(ip, port))
}
