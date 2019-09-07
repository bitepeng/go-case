package main

import "fmt"

func main() {
	d, err := Whois("whois.arin.net")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)

}
