# http-server-minimize
最小化的http服务器

## 核心代码
```
// Server
http.HandleFunc("/",handleHello)
if err:=http.ListenAndServe(":8080",nil);err!=nil{
    fmt.Println(err)
}

// Hello
func handleHello(w http.ResponseWriter,r *http.Request)  {
	if _,err:=w.Write([]byte("Hello World"));err!=nil{
		fmt.Println(err)
	}
}
```

## 使用说明
```
# 编译或运行
go build main.go
go run main.go
```
> 浏览器访问
http://127.0.0.1:8080