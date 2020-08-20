# A Unix Socket client for strongswan using [govici](https://github.com/strongswan/govici)


cat ~/export-http-proxy.sh
```sh
#!/bin/bash
export http_proxy=socks5://127.0.0.1:1080
export https_proxy=$http_proxy
export | grep http
```

# build

```
go install github.com/lixit/go_vici
export PATH=$PATH:$(dirname $(go list -f '{{.Target}}' .))
go_vici
```
