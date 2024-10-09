package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"net"
	"net/http"
	"strings"
)

func HttpRequestGetBody(req *http.Request) ([]byte, error) {
	var (
		reader  io.ReadCloser
		err     error
		reqBody []byte
	)
	switch req.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(req.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = req.Body
	}
	defer reader.Close()
	reqBody, err = io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return reqBody, nil
}

// ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
	start net.IP
	end   net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

var privateRanges = []ipRange{
	ipRange{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	ipRange{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	ipRange{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	ipRange{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	ipRange{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
	ipRange{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

// isPrivateSubnet - check to see if this ip is in a private subnet
func isPrivateSubnet(ipAddress net.IP) bool {
	// my use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range privateRanges {
			// check if this ip is in a private range
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

var httpProxyIPReversed = false

func SetHttpClientIPAddressBeReversed(reversed bool) {
	httpProxyIPReversed = reversed
}

func GetHttpClientIPAddress(r *http.Request) string {
	return getHttpClientIPAddressImpl(r, httpProxyIPReversed)
}

func getHttpClientIPAddressImpl(r *http.Request, beReversed bool) string {
	remoteAddr := r.RemoteAddr
	for _, h := range []string{"X-Forwarded-For", "x-forwarded-for", "X-Real-Ip", "X-Real-IP"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		if httpProxyIPReversed {
			for i := 0; i < len(addresses); i++ {
				ip := strings.TrimSpace(addresses[i])
				// header can contain spaces too, strip those out.
				realIP := net.ParseIP(ip)
				if !realIP.IsGlobalUnicast() /*|| isPrivateSubnet(realIP)*/ {
					// bad address, go to next
					continue
				}
				return ip
			}
		} else {
			for i := len(addresses) - 1; i >= 0; i-- {
				ip := strings.TrimSpace(addresses[i])
				// header can contain spaces too, strip those out.
				realIP := net.ParseIP(ip)
				if !realIP.IsGlobalUnicast() /*|| isPrivateSubnet(realIP)*/ {
					// bad address, go to next
					continue
				}
				return ip
			}
		}
	}
	if remoteAddr == "::1" {
		return "127.0.0.1"
	}
	splitString := strings.Split(remoteAddr, ":")
	return splitString[0]
}
