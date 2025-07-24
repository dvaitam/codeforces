package main

import (
    "bufio"
    "fmt"
    "os"
)

type Edge struct{to int; w int}

func solve(reader *bufio.Reader, writer *bufio.Writer){
    var n,k int
    if _, err := fmt.Fscan(reader, &n, &k); err != nil {
        return
    }
    var m int
    fmt.Fscan(reader, &m)
    adj := make([][]Edge, n+1)
    for i:=0;i<m;i++{
        var u,v,c int
        fmt.Fscan(reader,&u,&v,&c)
        adj[u]=append(adj[u],Edge{v,c})
    }
    // state encoding
    Np := n+2
    encode := func(L,R,pos int) int {
        return ((L*Np+R)*n + (pos-1))
    }
    decode := func(code int)(L,R,pos int){
        pos = code% n +1
        code/=n
        R = code%Np
        L = code/Np
        return
    }
    const INF int64 = 1<<60
    dpPrev := map[int]int64{}
    for i:=1;i<=n;i++{
        idx := encode(0,n+1,i)
        dpPrev[idx]=0
    }
    for step:=1; step<k; step++{
        dpNext := map[int]int64{}
        for idx,cost := range dpPrev {
            L,R,pos := decode(idx)
            for _,e := range adj[pos]{
                next := e.to
                if next<=L || next>=R || next==pos {continue}
                var newL,newR int
                if next>pos {
                    newL = pos
                    newR = R
                }else{
                    newL = L
                    newR = pos
                }
                nidx := encode(newL,newR,next)
                nc := cost + int64(e.w)
                if old,ok := dpNext[nidx]; !ok || nc<old{
                    dpNext[nidx]=nc
                }
            }
        }
        dpPrev = dpNext
        if len(dpPrev)==0 {break}
    }
    ans := INF
    for _,cost := range dpPrev{
        if cost<ans{ ans=cost }
    }
    if ans==INF{ fmt.Fprintln(writer,-1) } else { fmt.Fprintln(writer,ans) }
}

func main(){
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    solve(reader,writer)
    writer.Flush()
}
