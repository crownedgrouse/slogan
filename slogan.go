package slogan

import (
	"io"
	"runtime"
	"log"
	"os"
 	"fmt"
 	"syscall"
 	"unsafe"
 	"path"
 	"time"
	"github.com/bclicn/color"    // colorize output
)

/*
 * NOTE : Some exported functions have another version with an underscore ending name
 *        This avoid having 'declared and not used' build errors when commenting traces on variables.
 *        Instead add an underscore to silent the trace, and can be removed later to reactivate the trace.
 */

var logger = log.New(os.Stderr, "", 0)

var start = time.Now()

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

var formats = map[string]string{
	"fatal"   : "Immediate exit with code %d", // immediate exit on error format
	"trace"   : "%[1]T\n%%v: %[1]v\n\n%%v+: %+[1]v\n\n%%#v: %#[1]v",
	"empty"   : "%#v",
	"runtime" : "OS:%s ARCH:%s CPU:%d COMPILER:%s ROOT:%s",
	"default" : "   %[1]s %[2]s",
	"caller"  : "   %[1]s %[3]s\t %[2]s",
	"where"   : "%s:%d",
	"elapsed" : "All done in : %s",
}

var colors = map[int]string{
	10: "Underline",
	9:  "DarkGray",
	8:  "DarkGray",
	7:  "Purple",
	6:  "Green",  
	5:  "Yellow", 
	4:  "LightRed",   
	3:  "Red",   
	2:  "BLightRed",    
	1:  "BRed",
	0:  "",
}

// What parts of log should be colorized if Colorize=true
var parts = map[string]bool{
	"caller": true,
	"tag":    true,
	"log":    false,
	"prefix": false,
}

var offset = 0

var Verbosity   	int  = 5
var ExitOnError 	bool = false
var WarningAsError 	bool = false
var TraceCaller 	bool = false
var CallerBase  	bool = true
var Colorize    	bool = true

//************ Exported functions for configuration *************

/* Set global verbosity */
func SetVerbosity(level int) {
	Verbosity = level
}

/* Set exit on level error or higher */
func SetExitOnError(mode bool) {
	ExitOnError = mode
}

func SetWarningAsError(){
	WarningAsError = true
}

/* Set caller information in Trace */
func SetTraceCaller(mode bool) {
	TraceCaller = mode
}

/* Colorize or not */
func SetColor(mode bool) {
	Colorize = mode
}

func GetColors()map[int]string{
	return colors
}

func ShowColors(){
	fmt.Printf("%#v\n", colors)
}

func SetColors(n map[int]string)map[int]string{
	old := colors
	colors = n
	return old
}

/* API for logger override */
func SetFlags(flag int){
	if (flag & log.Lshortfile) == log.Lshortfile {
		TraceCaller = true
		CallerBase  = true
		SetFlags(flag - log.Lshortfile)
	} else 	if (flag & log.Llongfile) == log.Llongfile{
		TraceCaller = true
		CallerBase  = false
		SetFlags(flag - log.Llongfile)
	} else {
		logger.SetFlags(flag)
	}
}

func SetPrefix(prefix string) string{
	defer logger.SetPrefix(prefix)
	old := tags[0]
	tags[0] = prefix
	return old
}


func SetOutput(w io.Writer){
	logger.SetOutput(w)
}

func AllDone(){
	elapsed := time.Since(start)
	defer resetStart()
    Notice(fmt.Sprintf(formats["elapsed"], elapsed))
}

func resetStart(){
	start = time.Now()
}

//*** Levels ***
func GetTags()[10]string{
	return tags
}

func ShowTags(){
	fmt.Printf("%#v\n", tags)
}

func SetTags(n [10]string)[10]string{
	old := tags
	tags = n
	return old
}
//*** Formats ***
func GetFormats()map[string]string{
	return formats
}

func ShowFormats(){
	fmt.Printf("%#v\n", formats)
}

func SetFormats(n map[string]string)map[string]string{
	old := formats
	formats = n
	return old
}

//*** Parts ***
func GetParts()map[string]bool{
	return parts
}

func ShowParts(){
	fmt.Printf("%#v\n", parts)
}

func SetParts(n map[string]bool)map[string]bool{
	old := parts
	parts = n
	return old
}

//********** Exported functions for logging ****************************

// silent 0 | emergency 1 | alert 2 | critical 3 | error 4 | warning 5 | notice 6 | info 7 | debug 8 | trace 9

// Silent a log while keeping it
func Silent(log string) {
	Log(-1, log)
}

func Emergency(log string) {
	Log(1, log)
}

func Alert(log string) {
	Log(2, log)
}

func Critical(log string) {
	Log(3, log)
}

func Error(log string) {
	Log(4, log)
}

func Warning(log string) {
	Log(5, log)
}

func Notice(log string) {
	Log(6, log)
}

func Info(log string) {
	Log(7, log)
}

func Debug(log string) {
	Log(8, log)
}

// Trace non string as debug
func Trace(trace interface{}) {
	if fmt.Sprintf("%v", trace) == "[]" {
		Log(9, fmt.Sprintf(formats["empty"], trace))
	} else {
		Log(9, fmt.Sprintf(formats["trace"], trace))
	}	
}
func Trace_(trace interface{}){} 

// Trace non string as debug with caller punctually
func TraceCall(trace interface{}) {
	TraceCaller = true
	defer SetTraceCaller(false)
	Trace(trace)
}
func TraceCall_(trace interface{}){} 

func Runtime() {
	incr_offset()
	defer decr_offset()
    Debug(fmt.Sprintf(formats["runtime"], runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.Compiler, runtime.GOROOT()))
}

func incr_offset(){
	offset = offset + 1
}

func decr_offset(){
	offset = offset - 1
}

func Log(level int, log string) {

	if ( Verbosity >= level) {
		Str := logfmt(level, log)
		logger.Println(Str)
	}
	if (level < 5) && (ExitOnError == true) {
			Debug(fmt.Sprintf(formats["fatal"], level))
			os.Exit(level)
	}	
}

//****** Internal functions *************************************

func logfmt (level int, log string) string {
	Fmt := formats["default"]
	Tag := tags[level]

	Str := ""
	Caller := ""
	var fn_ string = ""
	var line int

	_, fn_, line, _ = runtime.Caller(3 + offset)

	if TraceCaller == true {
			fn := ""
			if CallerBase == true {
				fn = path.Base(fn_)
			}else {
				fn = fn_
			}
			Caller := colorize("caller", 10, fmt.Sprintf(formats["where"],fn, line))
			Str     = fmt.Sprintf(formats["caller"], colorize("tag", level,  Tag), colorize("log", level, log), Caller)
	}else{
		Str = fmt.Sprintf(Fmt, colorize("tag", level, Tag), colorize("log", level, log), Caller)	
	}
	return Str
}

func colorize(what string, level int, str string) string{
	if Colorize == true && parts[what] == true {
		return setcolor(what, level, str)
	}else{
		return str
	}
}

func setcolor(what string, level int, str string) string {
	Ret := ""
	switch colors[level] {
		case "Black": 		 Ret = color.Black(str)
		case "Red": 		 Ret = color.Red(str)
		case "Green": 		 Ret = color.Green(str)
		case "Yellow": 		 Ret = color.Yellow(str)
		case "Blue": 		 Ret = color.Blue(str)
		case "Purple": 		 Ret = color.Purple(str)
		case "Cyan": 		 Ret = color.Cyan(str)
		case "LightGray": 	 Ret = color.LightGray(str)
		case "DarkGray": 	 Ret = color.DarkGray(str)
		case "LightRed": 	 Ret = color.LightRed(str)
		case "LightGreen": 	 Ret = color.LightGreen(str)
		case "LightYellow":  Ret = color.LightYellow(str)
		case "LightBlue": 	 Ret = color.LightBlue(str)
		case "LightPurple":  Ret = color.LightPurple(str)
		case "LightCyan": 	 Ret = color.LightCyan(str)
		case "White": 		 Ret = color.White(str)
		// bold
		case "BBlack": 		 Ret = color.BBlack(str)
		case "BRed": 		 Ret = color.BRed(str)
		case "BGreen": 		 Ret = color.BGreen(str)
		case "BYellow": 	 Ret = color.BYellow(str)
		case "BBlue": 		 Ret = color.BBlue(str)
		case "BPurple": 	 Ret = color.BPurple(str)
		case "BCyan": 		 Ret = color.BCyan(str)
		case "BLightGray": 	 Ret = color.BLightGray(str)
		case "BDarkGray": 	 Ret = color.BDarkGray(str)
		case "BLightRed": 	 Ret = color.BLightRed(str)
		case "BLightGreen":  Ret = color.BLightGreen(str)
		case "BLightYellow": Ret = color.BLightYellow(str)
		case "BLightBlue": 	 Ret = color.BLightBlue(str)
		case "BLightPurple": Ret = color.BLightPurple(str)
		case "BLightCyan": 	 Ret = color.BLightCyan(str)
		case "BWhite": 		 Ret = color.BWhite(str)
		// background
		case "GBlack": 		 Ret = color.GBlack(str)
		case "GRed": 		 Ret = color.GRed(str)
		case "GGreen": 		 Ret = color.GGreen(str)
		case "GYellow": 	 Ret = color.GYellow(str)
		case "GBlue": 		 Ret = color.GBlue(str)
		case "GPurple": 	 Ret = color.GPurple(str)
		case "GCyan": 		 Ret = color.GCyan(str)
		case "GLightGray": 	 Ret = color.GLightGray(str)
		case "GDarkGray": 	 Ret = color.GDarkGray(str)
		case "GLightRed": 	 Ret = color.GLightRed(str)
		case "GLightGreen":  Ret = color.GLightGreen(str)
		case "GLightYellow": Ret = color.GLightYellow(str)
		case "GLightBlue": 	 Ret = color.GLightBlue(str)
		case "GLightPurple": Ret = color.GLightPurple(str)
		case "GLightCyan": 	 Ret = color.GLightCyan(str)
		case "GWhite": 		 Ret = color.GWhite(str)
		// special
		case "Bold": 		 Ret = color.Bold(str) 
		case "Dim": 		 Ret = color.Dim(str) 
		case "Underline": 	 Ret = color.Underline(str)
		case "Invert": 		 Ret = color.Invert(str) 
		case "Hide": 		 Ret = color.Hide(str)
		case "Blink": 		 Ret = color.Blink(str) // blinking works only on mac
		default: 
			Ret = str
	}
	return Ret
}

/*
 *   Terminal
 */
type winsize struct {
    Row    uint16
    Col    uint16
    Xpixel uint16
    Ypixel uint16
}

func getWidth() uint {
    ws := &winsize{}
    retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
        uintptr(syscall.Stdin),
        uintptr(syscall.TIOCGWINSZ),
        uintptr(unsafe.Pointer(ws)))

    if int(retCode) == -1 {
        panic(errno)
    }
    return uint(ws.Col)
}