package main
import "fmt"
func main(){
  for n:=1;n<=5;n++{
    fmt.Println("n=",n)
    count:=0
    for mask:=0;mask<1<<uint(n);mask++{
      arr:=make([]int,n)
      for i:=0;i<n;i++{
        if mask>>uint(i)&1==1{arr[i]=1}else{arr[i]=0}
      }
      if f(arr)==1{
        fmt.Println(arr)
        count++
      }
    }
    fmt.Println("count=",count)
  }
}
func f(arr []int) int {
  n:=len(arr)
  for i:=1;i<n;i++{
    if arr[i]-arr[i-1]!=1{
      return 1
    }
  }
  return 0
}
