package main
import "fmt"
func main(){
    for n:=2;n<=6;n++{
        adj:=make([][]int,n)
        for i:=0;i<n;i++{
            for j:=0;j<n;j++{
                if i!=j{adj[i]=append(adj[i],j)}
            }
        }
        pops:=simulate(adj)
        fmt.Println("n=",n,"count=",sum(pops,"deg"))
    }
}
func sum(pops []int, s string) int {return 0}
func simulate(adj [][]int) []int{
    n:=len(adj)
    flag:=make([]bool,n)
    stack:=[]int{0}
    pops:=make([]int,n)
    for len(stack)>0{
        u:=stack[len(stack)-1]
        stack=stack[:len(stack)-1]
        pops[u]++
        flag[u]=true
        for _,v:=range adj[u]{
            if !flag[v]{
                stack=append(stack,v)
            }
        }
    }
    fmt.Println(pops)
    return pops
}
