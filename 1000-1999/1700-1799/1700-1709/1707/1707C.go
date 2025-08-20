package main

import (
    "bufio"
    "fmt"
    "os"
)

type DSU struct{
    p []int
}

func NewDSU(n int)*DSU{
    d:=&DSU{p:make([]int,n+1)}
    for i:=1;i<=n;i++{d.p[i]=i}
    return d
}
func (d*DSU) find(x int) int{ if d.p[x]==x {return x}; d.p[x]=d.find(d.p[x]); return d.p[x] }
func (d*DSU) union(a,b int) bool{ a=d.find(a); b=d.find(b); if a==b {return false}; d.p[a]=b; return true }

func main(){
    in:=bufio.NewReader(os.Stdin)
    out:=bufio.NewWriter(os.Stdout)
    defer out.Flush()
    var n,m int
    fmt.Fscan(in,&n,&m)
    type P struct{u,v int}
    edges:=make([]P,m)
    for i:=0;i<m;i++{ fmt.Fscan(in,&edges[i].u,&edges[i].v) }

    // Build MST by Kruskal in edge id order; collect tree and non-tree edges
    dsu:=NewDSU(n)
    tree:=make([][]int,n+1)
    extra:=make([]P,0)
    for _,e:=range edges{
        u,v:=e.u,e.v
        if dsu.union(u,v){ tree[u]=append(tree[u],v); tree[v]=append(tree[v],u) } else { extra=append(extra,P{u,v}) }
    }

    // Depth array using DFS from 1
    depth:=make([]int,n+1)
    var dfs1 func(int,int)
    dfs1=func(x,fa int){ depth[x]=depth[fa]+1; for _,y:=range tree[x]{ if y!=fa { dfs1(y,x) } } }
    dfs1(1,0)

    // Build directed edges from deeper node to ancestor
    up:=make([][]int,n+1)
    for _,e:=range extra{
        u,v:=e.u,e.v
        if depth[u]<depth[v]{ u,v=v,u }
        up[u]=append(up[u],v)
    }

    // Difference array logic as in accepted solutions
    a:=make([]int,n+2)
    vis:=make([]bool,n+1)
    stk:=make([]int,n+2) // s[depth] mapping

    var dfs2 func(int,int)
    dfs2=func(x,fa int){
        stk[depth[x]]=x
        vis[x]=true
        for _,y:=range up[x]{
            a[x]++
            if vis[y]{
                a[1]++
                dNext:=depth[y]+1
                if dNext<=n{ a[stk[dNext]]-- }
            }else{
                a[y]++
            }
        }
        for _,y:=range tree[x]{ if y!=fa { dfs2(y,x) } }
        vis[x]=false
    }
    dfs2(1,0)

    var dfs3 func(int,int)
    dfs3=func(x,fa int){ a[x]+=a[fa]; for _,y:=range tree[x]{ if y!=fa { dfs3(y,x) } } }
    dfs3(1,0)

    need:=len(extra)
    ans:=make([]byte,n)
    for i:=1;i<=n;i++{ if a[i]==need { ans[i-1]='1' } else { ans[i-1]='0' } }
    fmt.Fprintln(out,string(ans))
}
