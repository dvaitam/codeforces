package main
import "fmt"

func simulate(adj [][]int) int {
    n:=len(adj)
    flag:=make([]bool,n)
    stack:=[]int{0}
    counter:=0
    for len(stack)>0 {
        u:=stack[len(stack)-1]
        stack=stack[:len(stack)-1]
        flag[u]=true
        for _,v:=range adj[u]{
            counter++
            if !flag[v] {
                stack=append(stack,v)
            }
        }
    }
    return counter
}

func main(){
    for a:=1;a<=5;a++{
        for b:=1;b<=5;b++{
            n:=a+b
            adj:=make([][]int,n)
            for i:=0;i<a;i++{
                for j:=0;j<b;j++{
                    u:=i
                    v:=a+j
                    adj[u]=append(adj[u],v)
                    adj[v]=append(adj[v],u)
                }
            }
            fmt.Println("K",a,b,"=",simulate(adj))
        }
    }
}
