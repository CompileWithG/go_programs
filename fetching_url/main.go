//To illustrate the minimum necessary to retreive information over HTTP,
//here's a simple program called fetch that fetches the content of each specified URL and priint it as a uninterpreted text;
//its inspired by the invaluable utility URL
package main 
import (

  "fmt"
  "io/ioutil"
  "net/http"
  "os"
)
//THe http.Get function makes an HTTP request and if there is no error ,returns the result in the response struct resp
func main(){
  for _, url := range os.Args[1:]{
    resp,err:= http.Get(url)
    if err!=nil{
      fmt.Fprintf(os.Stderr,"fetch: %v\n",err)
      os.Exit(1)
    }
    b,err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err!=nil{
      fmt.Fprintf(os.Stderr,"fetch: reading %s: %v\n",url,err)
      os.Exit(1)
    }
    fmt.Printf("%s",b)
  }
}
