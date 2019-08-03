# http-server-minimize
最小化的http服务器

```
http.Handle("/",http.FileServer(http.Dir(".")))
if err:=http.ListenAndServe(":8080",nil);err!=nil{
    fmt.Println(err)
}
```