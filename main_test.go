package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestDelRe(t *testing.T) {
	b, e := ioutil.ReadFile("test2.html")
	if e != nil {
		log.Fatal(e)
	}
	find := del_re.FindAllString(string(b), -1)
	log.Println(find[0])
	log.Println("wr_id =", del_re.FindStringSubmatch(find[0])[1])
}
