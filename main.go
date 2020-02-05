package main

import (
	"log"
  "fmt"
	"net"
  "net/http"
	"time"
	"bufio"
  "os"
	"regexp"
  "strconv"
  "strings"
)

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
	// fmt.Printf("match size: %v\n", len(match))
	defer conn.Close()
	if len(match) != 1 {
		return nil
	}
	return fmt.Errorf("%v", match[0])
}

func Reverse(s string) string {
    n := len(s)
    runes := make([]rune, n)
    for _, rune := range s {
        n--
        runes[n] = rune
    }
    return string(runes[n:])
}

func main() {
    
	  var PORT string

    if os.Getenv("PORT") != "" {
		  PORT = os.Getenv("PORT")
		} else {
		  PORT = "80"
    }

    http.HandleFunc("/address", AddressLookup)
    http.HandleFunc("/uuid/", UUIDLookup)

    fmt.Printf("Starting server at :%s...\n",PORT)

    log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func AddressLookup( w http.ResponseWriter, r *http.Request) {
    res := strings.SplitN(Reverse(r.RemoteAddr),":", -1)
    remote_ip := Reverse(strings.Join(res[1:], ":"))
    remote_ip = strings.Replace(remote_ip, "[", "", -1)
    remote_ip = strings.Replace(remote_ip, "]", "", -1)
    fmt.Fprintf(w, "%v\n", remote_ip)
}

func UUIDLookup(w http.ResponseWriter, r *http.Request) {
    res := strings.SplitN(Reverse(r.RemoteAddr),":", -1)
    remote_ip := Reverse(strings.Join(res[1:], ":"))
    remote_ip = strings.Replace(remote_ip, "[", "", -1)
    remote_ip = strings.Replace(remote_ip, "]", "", -1)

    ip  := net.ParseIP(remote_ip)

    path := strings.SplitN(string(r.URL.Path[1:]), "/", -1)
    port, _ := strconv.ParseUint(path[1],10,64)
    fmt.Fprintf(w, "%v\n", LookupUUID(ip, port))
}
