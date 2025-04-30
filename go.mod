module github.com/ab36245/go-msgpack

go 1.24.2

replace github.com/ab36245/go-errors => ../go-errors

require (
	github.com/ab36245/go-errors v0.0.0-20250428061939-8b056c3b905e
	github.com/vmihailenco/msgpack/v5 v5.4.1
)

require github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
