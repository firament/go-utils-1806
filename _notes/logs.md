# Logging in Go Lang

## To Read
- https://logmatic.io/blog/our-guide-to-a-golang-logs-world/
`Structuring your Golang logs in one project or across multiples microservices is probably the hardest part of the journey, even though it could seem trivial once done.`
- https://blog.gopheracademy.com/advent-2016/context-logging/
- Logrus/Logmatic.io
- github.com/op/go-logging
- github.com/cryptix/exp/wslog
- Rsyslog
- Logstash
- FluentD
- https://cloud.google.com/appengine/docs/standard/go/logs/
- https://docs.gocd.org/current/advanced_usage/logging.html

## golang.org/pkg/log/
- `func Output(calldepth int, s string) error`
    - Output writes the output for a logging event. 
    - The string s contains the text to print after the prefix specified by the flags of the Logger. 
    - A newline is appended if the last character of s is not already a newline. 
    - Calldepth is the count of the number of frames to skip when computing the file name and line number if Llongfile or Lshortfile is set; a value of 1 will print the details for the caller of Output. 

- `func SetFlags(flag int)`
    - SetFlags sets the output flags for the standard logger. 

- `func SetOutput(w io.Writer)`
    - SetOutput sets the output destination for the standard logger. 

## golang.org/pkg/log/syslog/
- Package syslog provides a simple interface to the system log service. It can send messages to the syslog daemon using UNIX domain sockets, UDP or TCP.
- Only one call to Dial is necessary. On write failures, the syslog client will attempt to reconnect to the server and write again. 

### See https://godoc.org/?q=syslog

## github.com/sirupsen/logrus

## github.com/golang/glog

## github.com/cihub/seelog

## github.com/uber-go/zap

