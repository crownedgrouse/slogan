# slogan #

`slogan` is a logger library for Golang.

Features :
   - 10 Levels from "silent" to "trace"
   - Configurable format
   - Configurable colors on individual items
   - Easy elapsed time

## Configuring ##

### Formats ###

Formats can be configured by settings new values

1- tag
2- log
3- caller (path:line)

```go
// Get current map
formats := slogan.GetFormats()
// Override format for trace : set caller in first
formats["trace"]= "%-25[3]s %[1]s %[2]s "
slogan.SetFormats(formats)
```


