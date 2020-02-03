# slogan #

`slogan` is a logger library for Golang.

Features :
   - 10 Levels from "silent" to "trace", including POSIX levels
   - Configurable format
   - Configurable colors on individual items
   - Easy elapsed time

## Howto ##

### Declare use ###

```go
import (
	"github.com/crownedgrouse/slogan"
)
```
An alias (here 'log') can be used by setting it before import :

```go
import (
	log  "github.com/crownedgrouse/slogan"
)
```

### Init ###

```go main.go
package main

import (
	log  	"github.com/crownedgrouse/slogan"
)

func main() {	
	
	log.SetExitOnError(true) 

	// Set verbosity
	log.SetVerbosity(9)
	//
	log.Runtime()
	log.Trace(map[string]string{"this is a": "map"})
	log.Debug("A debug message")
	log.Info("An informative message")
	log.Notice("A notification")
	log.Warning("A warning")
	log.Error("An Error")
	log.Critical("A critical message")
	log.Alert("An alert")
	log.Emergency("An Emergency")
}

```
Will produce (color not visible in this example):

```shell
$ go build
$ ./main
   debug     OS:linux ARCH:386 CPU:4 COMPILER:gc ROOT:/home/eric/git/goroot
   trace     map[string]string
%v: map[this is a:map]

%v+: map[this is a:map]

%#v: map[string]string{"this is a":"map"}
   debug     A debug message
   info      An informative message
   notice    A notification
   warning   A warning
   error     An Error
   debug     Immediate exit with code 4
$ echo $?
4
```

## Utilities ##

### Show Runtime infos ###

Runtime informations can be easily shown as a debug level with a uniq call.

```go
slogan.Runtime()
```
```shell
   debug     OS:linux ARCH:386 CPU:4 COMPILER:gc ROOT:/home/eric/git/goroot
```

### Trace Go values ###

Call to `Trace/1` will produce a trace log made of several lines. First line with 'trace' level and type of the value given. Below is written three usual ways to display Go values (%v, %v+ and %#v) separated with an empty line.


```go
	slogan.Trace(Something)
```

### Time elapsed ###

Display how many time elapsed since program start or since last call to `AllDone()` .

```go
    // Show time elapsed since beginning
    slogan.AllDone()
    // Show time elapsed since last call to AllDone()
    slogan.AllDone()
```
Output will be a notice :

```go
  notice    All done in : 296.538µs
```

## Configuring ##

`slogan` can be configured at beginning of your program (and also at any time inside your program).

### Tags ###

Tags can be changed by overwritting `tags` map, with `GetTags/0` and `SetTags/1`.

```go
var tags = [10]string{
		"",          // Prefix
		"emergency", // 1
		"alert    ", // 2
		"critical ", // 3
		"error    ", // 4
		"warning  ", // 5
		"notice   ", // 6
		"info     ", // 7
		"debug    ", // 8
		"trace    ", // 9
}

```

### Output ###

`slogan` is using legacy "log" package underneath. `SetFlags` can be used to change "log" parameters.

For instance to show caller and line number in code :

```go
	slogan.SetFlags(log.Lshortfile)  // Need to import "log" for this
```

Will produce something like below : 

```shell
   debug     main.go:17     A debug message
   info      main.go:18     An informative message
   notice    main.go:20     A warning
   error     main.go:21     An Error
```
as well date/time information can be set this way.

Set a prefix to any log :

```go
	log.SetPrefix("===> ")
```
### Behaviour ###

Considering Warning as Error (and potentialy exit) :

```go
	log.SetWarningAsError()
```
Set or change verbosity from 0 (silent) to 9 (debug) :

```go
slogan.SetVerbosity(0) // Silent totally logs
```
By setting verbosity, all logs with level lower or equal will be generated if no immediate exit on error was set :

```go
log.SetExitOnError(true) // Exit if log level reach Error or higher.
```
If the case, the error message is generated and a debug level may appear, depending current verbosity, indicating that an immediate exit occured, and telling what it the program exit code. The exit code is equal to the level reached by the last fatal error, i.e 1 (emergency) to 4 (error) , or even 5 if warning considered error.

### Formats ###

Formats can be configured by settings new "Sprintf" values to the three arguments passed to `slogan` functions :

- tag                  %[1]s
- log                  %[2]s
- caller (path:line)   %[3]s  (only if caller is required in "log" parameters)

```go
// Get current map
formats := slogan.GetFormats()
// Override format for caller : set caller in first
formats["caller"]= "%-25[3]s %[1]s %[2]s "
slogan.SetFormats(formats)
slogan.SetFlags(log.Lshortfile) // caller is required to be shown
```

Default formats are : 
```go
var formats = map[string]string{
	"fatal"   : "Immediate exit with code %d",                       // immediate exit on error format
	"trace"   : "%[1]T\n%%v: %[1]v\n\n%%v+: %+[1]v\n\n%%#v: %#[1]v", // multiline trace format
	"empty"   : "%#v",                                               // trace format for empty variable (avoid unuseful multiline)
	"runtime" : "OS:%s ARCH:%s CPU:%d COMPILER:%s ROOT:%s",          // runtime infos format
	"default" : "   %[1]s %[2]s",                                    // default log format
	"caller"  : "   %[1]s %[3]s\t %[2]s",                            // default log format with caller (where)
	"where"   : "%s:%d",                                             // format for caller location path:linenumber
	"elapsed" : "All done in : %s",                                  // elapsed time format
}
``` 

### Colors ###

Colors can be changed by overwritting `colors` map, with `GetColors/0` and `SetColors/1`.

See [here](https://github.com/bclicn/color) for possible colors and other output (reverse, underlining, etc.)

```go
var colors = map[int]string{
	10: "Underline",    // Caller
	9:  "DarkGray",     // trace
	8:  "DarkGray",     // debug
	7:  "Purple",       // info
	6:  "Green",        // notice
	5:  "Yellow",       // warning
	4:  "LightRed",     // error
	3:  "Red",          // critical
	2:  "BLightRed",    // alert
	1:  "BRed",         // emergency
	0:  "",             // Silent
}
```

As well colorization of elements (called 'parts') in log line can be tuned by changing `parts` map, with `GetParts/0` and `SetParts/1`

```go
var parts = map[string]bool{
	"caller": true,            // colorize caller (event if it is underlining)
	"tag":    true,            // colorize tag
	"log":    false,           // do not colorize log entry
	"prefix": false,           // do not colorize prefix
}
```



