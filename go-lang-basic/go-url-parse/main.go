package main

import (
	"fmt"
	"net/url"
	"strings"
)

func main() {

	// 我们将解析这个URL，它包含了模式，验证信息，
	// 主机，端口，路径，查询参数和查询片段
	urls := "mysql://user:pass@host.com:5432/path?key=value#f"

	u, _ := url.Parse(urls)
	uP, _ := u.User.Password()
	uH := strings.Split(u.Host, ":")
	uM, _ := url.ParseQuery(u.RawQuery)

	fmt.Println(
		urls,
		"|\n", // mysql
		u.Scheme,
		"|\n", // user:pass user pass
		u.User,
		u.User.Username(),
		uP,
		"|\n", // host.com:5432 host.com 5432 5432
		u.Host,
		uH[0],
		uH[1],
		u.Port(),
		"|\n", // /path key=value map[key:[value]] value f
		u.Path,
		u.RawQuery,
		uM,
		uM["key"][0],
		u.Fragment,
	)

}
