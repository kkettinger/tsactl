module github.com/kkettinger/tsactl

go 1.24.1

replace github.com/alecthomas/kong => github.com/kkettinger/kong v1.10.1-short-tag-fix

require (
	github.com/alecthomas/kong v1.10.0
	github.com/govalues/decimal v0.1.36
	github.com/kkettinger/go-tinysa v0.3.0
)

require (
	github.com/creack/goselect v0.1.3 // indirect
	go.bug.st/serial v1.6.4 // indirect
	golang.org/x/sys v0.32.0 // indirect
)
