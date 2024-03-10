package utils

import (
	"log"
	"net"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Path       = filepath.Dir(filepath.Dir(filepath.Dir(b))) + "/"
)

func GetIP(r *http.Request) string {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String()
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Fatalln(err)
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1"
		}
		return ip
	}

	log.Fatalln(err)
	return ""
}

func GetCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func CheckEmail(email string) bool {
	reg := regexp.MustCompile(`^[\w&#$.%+-]+@[\w&#$.%+-]+\.[a-z]{2,6}?$`)
	return reg.MatchString(email)
}

// CheckPasswd
// checks if the password's format is according to the rules.
func CheckPasswd(passwd string) bool {

	// Matches any password containing at least one digit, one lowercase,
	// one uppercase, one symbol and 8 characters in total.
	//regex := regexp.MustCompile(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*([^\w\s]|_)).{8,}$`) // Alas not supported by the regexp library
	digit := regexp.MustCompile(`\d+`)
	lower := regexp.MustCompile(`[a-z]+`)
	upper := regexp.MustCompile(`[A-Z]+`)
	symbol := regexp.MustCompile(`([^\w\s]|_)+`)
	minLen := regexp.MustCompile(`.{8,}`)
	return digit.MatchString(passwd) && lower.MatchString(passwd) && upper.MatchString(passwd) && symbol.MatchString(passwd) && minLen.MatchString(passwd)
}
func CheckPseudo(pseudo string) bool {
	reg := regexp.MustCompile(`^[A-Za-z0-9]{5,}$`)
	return reg.MatchString(pseudo)
}
