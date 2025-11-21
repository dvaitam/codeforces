package main
import (
    "fmt"
    "math/rand"
)

func brute(a []int, b []int) int {
    memo = map[int]map[int]int{}
    return solve(a,b,0,1)
}

var memo map[int]map[int]int

func solve(a []int, b []int, mask int, cur int) int {
    if memo[mask] == nil {memo[mask]=map[int]int{}}
    if v,ok:=memo[mask][cur]; ok {return v}
    idx:=cur-1
    best:=0
    // submit
    nmask:=mask|(1<<idx)
    nxt:=nextIdx(nmask, cur-1)
    cand:=a[idx]
    if nxt!=0 {cand+=solve(a,b,nmask,nxt)}
    if cand>best {best=cand}
    // skip
    nxt=nextIdx(nmask, b[idx])
    cand=0
    if nxt!=0 {cand+=solve(a,b,nmask,nxt)}
    if cand>best {best=cand}
    memo[mask][cur]=best
    return best
}

func nextIdx(mask int,bound int) int {
    for i:=bound;i>=1;i-- {
        if mask&(1<<(i-1))==0 {return i}
    }
    return 0
}

func greedySubmitAll(a []int, b []int) int {
    mask:=0
    cur:=1
    score:=0
    for cur!=0 {
        idx:=cur-1
        // always submit
        score+=a[idx]
        mask|=1<<idx
        cur=nextIdx(mask, cur-1)
    }
    return score
}

func main(){
    for n:=1;n<=6;n++{
        trials:=1000
        for t:=0;t<trials;t++{
            a:=make([]int,n)
            b:=make([]int,n)
            for i:=0;i<n;i++{
                a[i]=rand.Intn(5)+1
                b[i]=rand.Intn(n)+1
            }
            opt:=brute(a,b)
            g:=greedySubmitAll(a,b)
            if opt<g {panic("opt<g")}
        }
        fmt.Println("n",n,"ok")
    }
    fmt.Println("done")
}
