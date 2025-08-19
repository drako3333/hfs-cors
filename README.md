This is a modified version of `net/http.FileServer` to satisfy my use cases.

1. features can be enabled/disabled independently
   * access log (IP & URL path)
   * directory listing
   * serving index.html
   * serving dotfiles
   * serving cors
2. directory listing enhanced
   * add ../ entry
   * customizable sorting
     * list directories before files
     * case insensitive sorting
     * whatever order you like

```
$ go get -u github.com/zhangyoufu/hfs/cmd/hfs
$ hfs -help
Usage of hfs:
  -addr address:port
        listen address:port (default "localhost:8000")
  -dirfirst
        list directories before files (default true)
  -dotfile
        enable listing and serving dot files (default false)
  -igncase
        case insensitive sorting (default true)
  -index
        enable serving index.html (default true)
  -list
        enable directory listing (default true)
  -log
        enable access log (default true)
  -root path
        path of document root (default ".")
```
