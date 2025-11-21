package main
import (
    "fmt"
)

var best int

func solve(n int, edges [][2]int) int {
    adj := make([][]int,n)
    for _,e:= range edges{a,b:= e[0],e[1];adj[a]=append(adj[a],b);adj[b]=append(adj[b],a)}
    used := make([][]bool,n)
    for i:=0;i<n;i++{used[i]=make([]bool,len(adj[i]))}
    // operations remove edges along path length>=2 comprised of currently 
    // "available" edges (not removed). We'll brute via recursion.
    best = 1<<30
    dfs := func(remEdges int){}
    var rec func(adj [][]int)
    rec = func(adj [][]int){
        if lenadj := func() int {
            s:=0
            for i:=0;i<n;i++{s+=len(adj[i])}
            return s/2
        }(); lenadj==0{
            if 0<best{best=0}
            return
        }
    }
    return 0
}

func main(){fmt.Println("skip")}
