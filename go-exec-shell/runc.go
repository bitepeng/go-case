package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func main() {

	/*
		in :=bytes.NewBuffer(nil)
		cmd:=exec.Command("testc")
		cmd.Stdin =in
		go func() {
			in.WriteString("1 2 3")
			in.WriteString("exit \n")
		}()
		if err:=cmd.Run();err!=nil{
			fmt.Println(err)
			return
		}*/

	cmdexe := "testc"
	cmdin := "11 22 333"

	//生成可执行文件
	if err := RunCommand("gcc", cmdexe+".c", "-o", cmdexe); err != nil {
		fmt.Println(err)
	}
	//执行代码
	cmd := exec.Command(cmdexe, "")
	cmd.Stdin = strings.NewReader(cmdin)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cmdexe, cmdin)
	fmt.Println(out.String())
	//删除可执行文件
	if err := RunCommand("rm", "-rf", "testc"); err != nil {
		fmt.Println(err)
	}

}

func RunCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	// 命令的错误输出和标准输出都连接到同一个管道
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}
	// 从管道中实时获取输出并打印到终端
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}
