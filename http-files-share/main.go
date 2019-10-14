/* Tiny web server in Golang for sharing a folder
Copyright (c) 2019 bitepeng <bitepeng@qq.com>

Contains some code from Golang's http.ServeFile method, and
uses lighttpd's directory listing HTML template. */

package main

import "net/http"
import "net/url"
import "io"
import "os"
import "mime"
import "path"
import "fmt"
import "flag"
import "strings"
import "strconv"
import "text/template"
import "container/list"
import "compress/gzip"
import "compress/zlib"
import "time"

var rootFolder *string
var usesGzip *bool

// Server User Agent
const serverUA = "GoFiles/0.2"

// 4096 bits = default page size on OSX
const fsMaxbufsize = 4096

func main() {
	// Get current working directory to get the file from it
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while getting current directory.")
		return
	}

	// Command line parsing
	bind := flag.String("bind", ":1718", "Bind address")
	rootFolder = flag.String("root", cwd, "Root folder")
	usesGzip = flag.Bool("gzip", true, "Enables gzip/zlib compression")

	flag.Parse()

	http.Handle("/", http.HandlerFunc(handleFile))

	fmt.Printf("Sharing %s on %s ...\n", *rootFolder, *bind)
	if err = http.ListenAndServe(*bind, nil); err != nil {
		fmt.Printf("Error while ListenAndServe.")
		return
	}
}

// Manages directory listings
type dirlisting struct {
	Name          string
	ChildrenDir   []string
	ChildrenFiles []string
	ServerUA      string
}

func copyToArray(src *list.List) []string {
	dst := make([]string, src.Len())

	i := 0
	for e := src.Front(); e != nil; e = e.Next() {
		dst[i] = e.Value.(string)
		i = i + 1
	}

	return dst
}

func handleDirectory(f *os.File, w http.ResponseWriter, req *http.Request) {
	names, _ := f.Readdir(-1)

	// First, check if there is any index in this folder.
	for _, val := range names {
		if val.Name() == "index.html" {
			serveFile(path.Join(f.Name(), "index.html"), w, req)
			return
		}
	}

	// Otherwise, generate folder content.
	childrenDirTmp := list.New()
	childrenFilesTmp := list.New()

	for _, val := range names {
		if val.Name()[0] == '.' {
			continue
		} // Remove hidden files from listing

		if val.IsDir() {
			childrenDirTmp.PushBack(val.Name())
		} else {
			childrenFilesTmp.PushBack(val.Name())
		}
	}

	// And transfer the content to the final array structure
	childrenDir := copyToArray(childrenDirTmp)
	childrenFiles := copyToArray(childrenFilesTmp)

	tpl, err := template.New("tpl").Parse(dirlistingTpl)
	if err != nil {
		http.Error(w, "500 Internal Error : Error while generating directory listing.", 500)
		fmt.Println(err)
		return
	}

	data := dirlisting{Name: req.URL.Path, ServerUA: serverUA,
		ChildrenDir: childrenDir, ChildrenFiles: childrenFiles}

	err = tpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

// serveFile 保存文件
func serveFile(filepath string, w http.ResponseWriter, req *http.Request) {
	// Opening the file handle
	f, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "404 Not Found : Error while opening the file.", 404)
		return
	}

	defer func() {
		_ = f.Close()
	}()

	// Checking if the opened handle is really a file
	statinfo, err := f.Stat()
	if err != nil {
		http.Error(w, "500 Internal Error : stat() failure.", 500)
		return
	}

	if statinfo.IsDir() { // If it's a directory, open it !
		handleDirectory(f, w, req)
		return
	}

	if (statinfo.Mode() &^ 07777) == os.ModeSocket { // If it's a socket, forbid it !
		http.Error(w, "403 Forbidden : you can't access this resource.", 403)
		return
	}

	// Manages If-Modified-Since and add Last-Modified (taken from Golang code)
	if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && statinfo.ModTime().Unix() <= t.Unix() {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Last-Modified", statinfo.ModTime().Format(http.TimeFormat))

	// Content-Type handling
	query, err := url.ParseQuery(req.URL.RawQuery)

	if err == nil && len(query["dl"]) > 0 { // The user explicitedly wanted to download the file (Dropbox style!)
		w.Header().Set("Content-Type", "application/octet-stream")
	} else {
		// Fetching file's mimetype and giving it to the browser
		if mimetype := mime.TypeByExtension(path.Ext(filepath)); mimetype != "" {
			w.Header().Set("Content-Type", mimetype)
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}
	}

	// Manage Content-Range (TODO: Manage end byte and multiple Content-Range)
	if req.Header.Get("Range") != "" {
		startByte := parseRange(req.Header.Get("Range"))

		if startByte < statinfo.Size() {
			_, _ = f.Seek(startByte, 0)
		} else {
			startByte = 0
		}

		w.Header().Set("Content-Range",
			fmt.Sprintf("bytes %d-%d/%d", startByte, statinfo.Size()-1, statinfo.Size()))
	}

	// Manage gzip/zlib compression
	outputWriter := w.(io.Writer)

	isCompressedReply := false

	if (*usesGzip) == true && req.Header.Get("Accept-Encoding") != "" {
		encodings := parseCSV(req.Header.Get("Accept-Encoding"))

		for _, val := range encodings {
			if val == "gzip" {
				w.Header().Set("Content-Encoding", "gzip")
				outputWriter = gzip.NewWriter(w)

				isCompressedReply = true

				break
			} else if val == "deflate" {
				w.Header().Set("Content-Encoding", "deflate")
				outputWriter = zlib.NewWriter(w)

				isCompressedReply = true

				break
			}
		}
	}

	if !isCompressedReply {
		// Add Content-Length
		w.Header().Set("Content-Length", strconv.FormatInt(statinfo.Size(), 10))
	}

	// Stream data out !
	buf := make([]byte, min(fsMaxbufsize, statinfo.Size()))
	n := 0
	for err == nil {
		n, err = f.Read(buf)
		_, _ = outputWriter.Write(buf[0:n])
	}

	// Closes current compressors
	switch outputWriter.(type) {
	case *gzip.Writer:
		_ = outputWriter.(*gzip.Writer).Close()
	case *zlib.Writer:
		_ = outputWriter.(*zlib.Writer).Close()
	}

	_ = f.Close()
}

func handleFile(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Server", serverUA)

	filepath := path.Join(*rootFolder, path.Clean(req.URL.Path))
	serveFile(filepath, w, req)

	fmt.Printf("\"%s %s %s\" \"%s\" \"%s\"\n",
		req.Method,
		req.URL.String(),
		req.Proto,
		req.Referer(),
		req.UserAgent()) // TODO: Improve this crappy logging
}

func parseCSV(data string) []string {
	splitted := strings.SplitN(data, ",", -1)

	data_tmp := make([]string, len(splitted))

	for i, val := range splitted {
		data_tmp[i] = strings.TrimSpace(val)
	}

	return data_tmp
}

func parseRange(data string) int64 {
	stop := (int64)(0)
	part := 0
	for i := 0; i < len(data) && part < 2; i = i + 1 {
		if part == 0 { // part = 0 <=> equal isn't met.
			if data[i] == '=' {
				part = 1
			}

			continue
		}

		if part == 1 { // part = 1 <=> we've met the equal, parse beginning
			if data[i] == ',' || data[i] == '-' {
				part = 2 // part = 2 <=> OK DUDE.
			} else {
				if 48 <= data[i] && data[i] <= 57 { // If it's a digit ...
					// ... convert the char to integer and add it!
					stop = (stop * 10) + (((int64)(data[i])) - 48)
				} else {
					part = 2 // Parsing error! No error needed : 0 = from start.
				}
			}
		}
	}

	return stop
}

func min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

const dirlistingTpl = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">
<head>
<title>Index of {{.Name}}</title>
<style type="text/css">
body{font-size:14px;line-height:180%;}
a, a:active {text-decoration: none; color: blue;}
a:visited {color: #48468F;}
a:hover, a:focus {text-decoration: underline; color: red;}
body {background-color: #F5F5F5;}
h2 {margin-bottom: 12px;}
table {margin-left: 12px;}
th, td { font: 90% monospace; text-align: left;}
th { font-weight: bold; padding-right: 14px; padding-bottom: 3px;}
td {padding-right: 14px;}
td.s, th.s {text-align: right;}
div.list { background-color: white; border-top: 1px solid #646464; border-bottom: 1px solid #646464; padding-top: 10px; padding-bottom: 14px;}
div.foot { font: 90% monospace; color: #787878; padding-top: 4px;}
</style>
</head>
<body>
<h2>Index of {{.Name}}</h2>
<div class="list">
<table summary="Directory Listing" cellpadding="0" cellspacing="0">
<thead><tr><th class="n">Name</th><th class="t">Type</th><th class="dl">Options</th></tr></thead>
<tbody>
<tr><td class="n"><a href="../">Parent Directory</a>/</td><td class="t">Directory</td><td class="dl"></td></tr>
{{range .ChildrenDir}}
<tr><td class="n"><a href="{{.}}/">{{.}}/</a></td><td class="t">Directory</td><td class="dl"></td></tr>
{{end}}
{{range .ChildrenFiles}}
<tr><td class="n"><a href="{{.}}">{{.}}</a></td><td class="t">&nbsp;</td><td class="dl"><a href="{{.}}?dl">Download</a></td></tr>
{{end}}
</tbody>
</table>
</div>
<div class="foot">{{.ServerUA}}</div>
</body>
</html>`
