package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jasonknight/gopress"
	"io/ioutil"
	"log"
	"os"
)

const (
	STARTUP = iota
)

func oops(s string) error {
	return errors.New(s)
}

var (
	sarg            = flag.String(`s`, `default value`, `document the option here`)
	logFilePath     = flag.String(`l`, `woo2ebay.log`, `the path to your chosen logfile`)
	yamlAdapterPath = flag.String(`a`, `../woo2ebay.db.yml`, `the adapter YAML for gopress`)
)
var Info *log.Logger
var Error *log.Logger
var mysql *gopress.MysqlAdapter

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
	Info = log.New(file, `[woo2ebay INF]:`, log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(file, `[woo2ebay ERR]:`, log.Ldate|log.Ltime|log.Lshortfile)
	Info.Println("Beginning")
	mysql, err = gopress.NewMysqlAdapterEx(*yamlAdapterPath)
	if err != nil {
		Error.Println(err)
		return
	}
	Info.Println("Database opened for reading")
	err = maybeSendProducts(mysql)
	if err != nil {
		Error.Println(err)
		return
	}
	return
}

func maybeSendProducts(a gopress.Adapter) error {
	Info.Println("Beginning maybeSendProducts")
	m := gopress.NewPostMeta(a)
	metas, err := m.FindByKeyValue("_woo2ebay_send", "yes")
	if err != nil {
		return err
	}
	if len(metas) == 0 {
		Info.Println("No metas found with send = yes")
		return nil
	}
	var products []*gopress.Post
	for _, meta := range metas {
		p := gopress.NewPost(a)
		found, err := p.Find(meta.PostId)
		if err != nil {
			return err
		}
		if found == true {
			products = append(products, p)
		}
	}
	if len(products) == 0 {
		Info.Println("could not find any products")
		return nil
	}
	return nil
}

// Helpers
func fileGetContents(p string) ([]byte, error) {
	return ioutil.ReadFile(p)
}

func fileExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}
