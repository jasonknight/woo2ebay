package main
import (
    "fmt"
    "errors"
    "flag"
    "log"
    "os"
)
func oops(s string) error {
    return errors.New(s)
}

var (
    sarg = flag.String(`s`, `default value`, `document the option here`)
    logFilePath = flag.String(`l`,`woo2ebay.log`,`the path to your chosen logfile`)
)
var Logger *log.Logger

func init() {
    flag.Parse()
}
func main() {
    file, err := os.OpenFile(*logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        fmt.Printf("Failed to open log file:", err)
        return
    }
    defer file.Close()
    Logger = log.New(file, `[woo2ebay]:`, log.Ldate|log.Ltime|log.Lshortfile)
    Logger.Println("Beginning")
}