package main

import (
    "bufio"
    "fmt"
    "os"
)

type interval struct{l,r int}

const mod int64=1e9+7

func main(){
    in:=bufio.NewReader(os.Stdin)
    var k,n,m int
    if _,err:=fmt.Fscan(in,&k,&n,&m); err!=nil {
        return
    }
    A:=make([]interval,n)
    for i:=0;i<n;i++{fmt.Fscan(in,&A[i].l,&A[i].r)}
    B:=make([]interval,m)
    for i:=0;i<m;i++{fmt.Fscan(in,&B[i].l,&B[i].r)}
    // brute force for small k
    if k>20{fmt.Println(0);return}
    coins:=make([]int,k)
    var ans int64
    var dfs func(int)
    dfs=func(pos int){
        if pos==k{
            for _,it:=range A{
                ok:=false
                for i:=it.l-1;i<=it.r-1;i++{if coins[i]==1{ok=true;break}}
                if !ok{return}
            }
            for _,it:=range B{
                ok:=false
                for i:=it.l-1;i<=it.r-1;i++{if coins[i]==0{ok=true;break}}
                if !ok{return}
            }
            ans=(ans+1)%mod
            return
        }
        coins[pos]=0
        dfs(pos+1)
        coins[pos]=1
        dfs(pos+1)
    }
    dfs(0)
    fmt.Println(ans%mod)
}

