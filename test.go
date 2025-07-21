package main
import "fmt"
func main() {
 n:=3
 size:=1; for size<n {size<<=1}
 x:=3
 l:=size; r:=size+x-1
 var nodes []int
 for l<=r {
  if l&1==1 { nodes=append(nodes,l); l++ }
  if r&1==0 { nodes=append(nodes,r); r-- }
  l>>=1; r>>=1
 }
 fmt.Println(size, nodes)
}
