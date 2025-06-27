# bdmv_go

It reads Blu-Ray information.



### Build mpls
```bash
$ go build ./pkg/mpls
```

### Build mpls-dump
```bash
$ go build -ldflags="-s -w" -o bin/mpls-dump ./cmd/mpls-dump
```