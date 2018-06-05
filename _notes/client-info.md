# Client Info

## libs
- github.com/tomasen/realip
    - `clientIP := realip.FromRequest(r)`
- github.com/rdegges/go-ipify
    - get your public IP address
- golang.org/pkg/net
    - `func ParseMAC(s string) (hw HardwareAddr, err error)`

    
## Notes:
- https://gist.github.com/rucuriousyet/ab2ab3dc1a339de612e162512be39283
- https://husobee.github.io/golang/ip-address/2015/12/17/remote-ip-go.html
```go
func getIPAdress(r *http.Request) string {
    for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
        addresses := strings.Split(r.Header.Get(h), ",")
        // march from right to left until we get a public address
        // that will be the address right before our proxy.
        for i := len(addresses) -1 ; i >= 0; i-- {
            ip := strings.TrimSpace(addresses[i])
            // header can contain spaces too, strip those out.
            realIP := net.ParseIP(ip)
            if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
                // bad address, go to next
                continue
            }
            return ip
        }
    }
    return ""
}
```