package main 
import (
  "fmt"
  "os"
  "io"
  "io/ioutil"
  "net/http"
  "time"
)

//a goroutine is a concurrent function execution .A channel is a commincation mechanism that allows one goroutine to pass values of a specified type
//to another goroutine .The function main runs in a goroutine and the go statement creates additional goroutines
func main(){
  start:= time.Now()
  ch:=make(chan string)
  for _,url := range os.Args[1:]{
    go fetch(url,ch)//stars a goroutine
  }
  for range os.Args[1:]{
    fmt.Println(<-ch)// receive from channel channel
  }
  fmt.Printf("%.2fs elapsed\n",time.Since(start).Seconds())
}

func fetch(url string,ch chan<-string){
  start:= time.Now()
  resp,err:=http.Get(url)
  if err!=nil{
    ch<-fmt.Sprint(err)//send to channel ch 
    return
  }
  nbytes,err:= io.Copy(ioutil.Discard,resp.Body)
  resp.Body.Close()//dont leak resources
  if err!=nil{
    ch<- fmt.Sprintf("while reading %s: %v",url,err)
    return
  }
  secs:=time.Since(start).Seconds()
  ch<- fmt.Sprintf("%.2fs %7d %s",secs,nbytes,url)
}
