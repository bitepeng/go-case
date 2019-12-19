package main

import "fmt"

/**
# Go range函数
range可以使用它来遍历数组，切片和字典

- 当用于遍历数组和切片的时候，range函数返回索引和元素 index,value
- 当用于遍历字典的时候，range函数返回字典的键和值 key,value
- 当用于遍历字符串时，返回Unicode代码点
*/

func main() {

	// 1. 遍历切片
	nums := []int{2, 3, 4}
	sum := 0
	// 只取值
	for _, num := range nums {
		sum += num
	}
	fmt.Println(sum)
	// 取索引和值
	for index, num := range nums {
		if num == 3 {
			fmt.Println("index:", index)
		}
	}

	// 2. 遍历字典
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs {
		fmt.Printf("%s -> %s\n", k, v)
	}

	// 3. 遍历字符串
	for i, c := range "go" {
		fmt.Println(i, c, string(c))
	}

}
