package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	/**
	File Server
	*/
	//1.实现读取文件handler
	fileHandle := http.FileServer(http.Dir("./video"))
	//2.注册handler，处理前缀
	http.Handle("/video/", http.StripPrefix("/video", fileHandle))
	/**
	HTTP Server
	*/
	//注册路由
	http.HandleFunc("/api/upload", uploadHandler)
	http.HandleFunc("/api/list", getFileHandler)
	//启动服务
	_ = http.ListenAndServe(":8090", nil)
}

/**
api1:上传视频文件
*/
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	//1.限制客户端上传视频文件的大小
	var maxFileSize int64 = 30 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//2.获取上传的文件
	file, fileHeader, err := r.FormFile("uploadFile")

	//3.检查文件类型
	ret := strings.HasSuffix(fileHeader.Filename, ".mp4")
	if ret == false {
		http.Error(w, "Not a mp4 file", http.StatusInternalServerError)
		return
	}

	//4.生成文件内容
	md5Byte := md5.Sum([]byte(fileHeader.Filename + time.Now().String()))
	md5Str := fmt.Sprintf("%x", md5Byte)
	newFileName := md5Str + ".mp4"

	//5.写入文件
	dst, err := os.Create("./video/" + newFileName)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

/**
Api2:获取视频文件列表
*/
func getFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")

	files, _ := filepath.Glob("video/*")
	var ret []string
	for _, file := range files {
		ret = append(ret, "http://"+r.Host+"/video/"+filepath.Base(file))
	}
	retJson, _ := json.Marshal(ret)
	w.Write(retJson)
	return

}
