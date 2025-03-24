package main 

import (
  "fmt" 
  "bufio"
  "os"
)
//bufio package-one of its most usefull feature is a type called Scanner that reads input and breaks it into lines or words
//;its often the easiest NewScanner
//to process input that comes naturally in lines
//The Scanner reads from the pograms standard input .Each call to input.Scan() reads the next line and removes the newline character
//from the end ;the result can be retrieved by calling input.Text().THe Scan function returns true if there is a lines and false
//when there is no more input
func main(){
  counts :=make(map[string] int)
  files:=os.Args[1:]
  if len(files)==0{
    countLines(os.Stdin,counts)
  }else{
    for _,arg:=range files{
      f,err :=os.Open(arg)
      if err!=nil{
        //print a message on the standard error stream using Fprintf
        fmt.Fprintf(os.Stderr,"dup2: %v\n",err)
        continue
      }
      countLines(f,counts)
      f.Close()
    }
  }
  for line,n := range counts{
    if n>1{
      fmt.Printf("%d\t%s\n",n,line)
    }
  }
}

func countLines(f *os.File,counts map[string]int){
  input:= bufio.NewScanner(f)
  for input.Scan(){
    counts[input.Text()]++
  }
  //NOte: ignoring potential errors from input error
}
