# cmd-app-calculator
比较完整的工程化计算器

## 使用说明
```
# 编译或运行
go build *.go
go run *.go

# 测试软件包
$ go test ./simplemath/
ok      ./simplemath       0.014s
 
```

```
# 命令行访问
$ ./calc 
USAGE: calc command [arguments] ...

The commands are:
        add     Addition pf two values.
        sqrt    Square root of a non-negative value.

$ ./calc add 888 903232
Result:  904120

$ ./calc sqrt 8
Result:  2
```
