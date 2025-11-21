package main
import "fmt"

func simulate(adj [][]int) int {
    n:=len(adj)
    flag:=make([]bool,n)
    stack:=[]int{0}
    counter:=0
    for len(stack)>0 {
        u:=stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        flag[u]=true
        for _,v:=range adj[u] {
            counter++
            if !flag[v] {
                stack=append(stack,v)
            }
        }
    }
    return counter
}

func main(){
    for r:=1;r<=10;r++{
        n:=r+1
        adj:=make([][]int,n)
        for i:=1;i<n;i++{
            adj[0]=append(adj[0],i)
            adj[i]=append(adj[i],0)
        }
        fmt.Println("star",r,"->",simulate(adj))
    }
}
