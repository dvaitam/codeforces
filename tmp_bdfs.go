package main
import "fmt"
func main(){
    for n:=2;n<=7;n++{
        adj:=make([][]int,n)
        for i:=0;i<n;i++{
            for j:=0;j<n;j++{
                if i!=j{
                    adj[i]=append(adj[i],j)
                }
            }
        }
        // sort adjacency ascending not needed because inherent order j increasing
        fmt.Println("clique",n,"->",simulate(adj))
    }
}
func simulate(adj [][]int) int{
    n:=len(adj)
    flag:=make([]bool,n)
    stack:=[]int{0}
    counter:=0
    for len(stack)>0 {
        u:=stack[len(stack)-1]
        stack=stack[:len(stack)-1]
        flag[u]=true
        for _,v:=range adj[u]{
            if v<=n { // placeholder
            }
            counter++
            if !flag[v]{
                stack=append(stack,v)
            }
        }
    }
    return counter
}
