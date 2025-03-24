package main 
import(
  "fmt"
  "net/http"
  "log"
)


func main(){
  http.HandleFunc("/",handler)//each request calls handler
  log.Fatal(http.ListenAndServe("localhost:8000",nil))
}

func handler(w http.ResponseWriter,r *http.Request){
  fmt.Fprintf(w,"URL.Path = %q\n", r.URL.Path)
}
//The main function connects a handler function to incoming URLs that begin with /, which is all URLs and start a server listening for 
//incoming request on port 8000 .A request is represented as a struct of type http.Request which contain a number of related fields .
//one of which is the URL of the incoming
