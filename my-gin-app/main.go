package main 


import "fmt" 

func main(){
  a:=5
  fmt.Println(a)
  b:= make(map[string]int)
  b["karthik"]=9
  fmt.Printf("i just created a map ,the keys is %v and the value is %v\n","karthik",b["karthik"])
  c:=[5]int{1,2,3,4,5}
  for i,elem := range c{
    fmt.Printf("this is the index %v and this is the element %v and this is the sum of the index plus the element %v\n",i,elem,add(i,elem) )
  }

}

func add(a int,b int ) int{
  return a+b
}
