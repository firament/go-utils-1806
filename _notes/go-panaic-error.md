# Recovering gracefully from panic

- Recover is useful only when called inside deferred functions.
- `debug.PrintStack()` will give the stack trace

## 1
```go
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
```

## 2
- For a real-world example of panic and recover, see the json package from the Go standard library. 
- It decodes JSON-encoded data with a set of recursive functions. 
- When malformed JSON is encountered, the parser calls panic to unwind the stack to the top-level function call, which recovers from the panic and returns an appropriate error value 
- See the 'error' and 'unmarshal' methods of the decodeState type in decode.go.

## 3
```go
func stackDump(err *error, f interface{}) {
	fname := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	_, file, line, _ := runtime.Caller(4) // this skips the first 4 that are called under log.Panic()
	if r := recover(); r != nil {
		fmt.Printf("%s (recover): %v\n", fname, r)
		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	} else if err != nil && *err != nil {
		fmt.Printf("%s : %v\n", fname, *err)
	}

	buf := make([]byte, 1<<10)
	runtime.Stack(buf, false)
	fmt.Println("==> stack trace: [PANIC:", file, line, fname+"]")
	fmt.Println(string(buf))
}
```

## 4
- Using "go-errors/errors"
```go
defer func() {
    if err := recover(); err != nil {
        fmt.Println(errors.Wrap(err, 2).ErrorStack())
    }
}()
```