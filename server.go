package main

import (
  "io";
  "os";
  "fmt";
  "time";
  "strings";
  "strconv";
  "net/http";
  "crypto/md5";
  "html/template";
)

var ImagesDir = http.Dir("images")
var ImagesDirHandler = http.FileServer(ImagesDir)

func GenerateFileName(filename string) (string){
  tim := time.Now().Format("20060102150405")
  splitsName := strings.Split(filename,".")
  var lenfilename = len(splitsName)
  if lenfilename == 1{
    return filename+"_"+tim
  } else {
    var newFileName = ""
    for i:=0; i<lenfilename-1;i++{
      newFileName += splitsName[i]
    }
    newFileName += "_"+tim+"."+splitsName[lenfilename-1]
    return newFileName
  }
}

// upload images
func UploadImagesHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("method:", r.Method)
  if r.Method == "GET" {
    curtime := time.Now().Unix()
    h := md5.New()
    io.WriteString(h, strconv.FormatInt(curtime, 10))
    token := fmt.Sprintf("%x", h.Sum(nil))
    fmt.Println(token)
    t, _ := template.ParseFiles("upload.html")
    t.Execute(w, token)
  } else {
    // opening file uploaded from client
    r.ParseMultipartForm(32 << 20)
    file, handler, err := r.FormFile("uploadfile")
    defer file.Close()
    if err != nil {
      fmt.Println(err)
      return
    }
    
    // Creating a new file to copy in server
    fmt.Fprintf(w, "%v", handler.Header)
    newFileName := GenerateFileName(handler.Filename)
    newfile, err := os.OpenFile("images/"+newFileName, os.O_WRONLY|os.O_CREATE, 0666)
    defer newfile.Close()
    if err != nil {
      fmt.Println(err)
      return
    }
    io.Copy(newfile, file)
  }
}

func main() {
  http.HandleFunc("/", UploadImagesHandler)
  http.Handle("/images/", http.StripPrefix("/images/",ImagesDirHandler))
  fmt.Println("Running on localhost:9090")
  http.ListenAndServe(":9090", nil)
}
