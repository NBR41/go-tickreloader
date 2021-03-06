package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/NBR41/gotickreloader"
)

func main() {

	var i = 10
	var cli = gotickreloader.NewClient(
		3*time.Second,
		func(v ...interface{}) (interface{}, error) {
			fmt.Println("reload")
			var value = v[0].(*int)
			*value = *value + v[1].(int)
			return *value, nil
		},
		&i, 3,
	)

	cli.StartTickReload()
	test(cli)
	cli.StopTickReload()

	cli = gotickreloader.NewClient(
		1*time.Second,
		func(v ...interface{}) (interface{}, error) {
			fmt.Println("reload")
			return []string{"foo", "bar"}, nil
		},
	)

	cli.StartTickReload()
	test(cli)
	cli.StopTickReload()

	cli = gotickreloader.NewClient(
		1*time.Second,
		func(v ...interface{}) (interface{}, error) {
			fmt.Println("!!!!! reload !!!!")
			return nil, errors.New("reload error")
		},
	)

	cli.StartTickReload()
	test(cli)
	cli.StopTickReload()
}

func test(cli *gotickreloader.Client) {
	var ch = make(chan bool)
	for j := 0; j < 10; j++ {
		go func(cli *gotickreloader.Client, j int, ch chan bool) {
			for k := 0; k < 10; k++ {
				v, err := cli.Get()
				fmt.Println("j", j, "k", k, "Get", v, err)
				time.Sleep(1 * time.Second)
			}
			if j == 9 {
				ch <- true
			}
		}(cli, j, ch)
	}
	<-ch
}
