package main
import "fmt"

func simulate(n int, edges [][2]int) (int, []int) {
    adj:=make([][]int,n)
    for _,e:=range edges {
        u,v:=e[0]-1,e[1]-1
        adj[u]=append(adj[u],v)
        adj[v]=append(adj[v],u)
    }
    for i:=0;i<n;i++{
        // sort adjacency ascending
        for j:=0;j<len(adj[i]);j++{
            for k:=j+1;k<len(adj[i]);k++{
                if adj[i][j]>adj[i][k]{adj[i][j],adj[i][k]=adj[i][k],adj[i][j]}
            }
        }
    }
    flag:=make([]bool,n)
    stack:=[]int{0}
    counter:=0
    pops:=make([]int,n)
    for len(stack)>0{
        u:=stack[len(stack)-1]
        stack=stack[:len(stack)-1]
        flag[u]=true
        pops[u]++
        for _,v:=range adj[u]{
            counter++
            if !flag[v]{
                stack=append(stack,v)
            }
        }
    }
    return counter,pops
}

func main(){
    n:=5
    edges:=[][2]int{}
    for i:=2;i<=n;i++{
        edges=append(edges,[2]int{1,i})
    }
    edges=append(edges,[2]int{2,3})
    edges=append(edges,[2]int{2,4})
    edges=append(edges,[2]int{3,5})
    fmt.Println(simulate(n,edges))
}
