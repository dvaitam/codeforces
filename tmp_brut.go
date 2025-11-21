package main
import (
    "fmt"
)

var memo map[int]map[int]int // mask -> cur -> best

func solve(a []int, b []int, mask int, cur int) int {
    if memo[mask] == nil {
        memo[mask] = map[int]int{}
    }
    if v, ok := memo[mask][cur]; ok {return v}
    idx := cur-1
    best := 0
    // submit
    nmask := mask | (1<<idx)
    bound := cur-1
    nxt := nextIdx(nmask, bound)
    cand := a[idx]
    if nxt != 0 {
        cand += solve(a,b,nmask,nxt)
    }
    if cand > best {best = cand}

    // skip
    bound = b[idx]
    nxt = nextIdx(nmask, bound)
    cand = 0
    if nxt != 0 {
        cand += solve(a,b,nmask,nxt)
    }
    if cand > best {best = cand}

    memo[mask][cur] = best
    return best
}

func nextIdx(mask int, bound int) int {
    for i:=bound; i>=1; i-- {
        if mask&(1<<(i-1))==0 {
            return i
        }
    }
    return 0
}

func brute(a []int, b []int) int {
    memo = map[int]map[int]int{}
    return solve(a,b,0,1)
}

func main(){
    a:=[]int{15,16}
    b:=[]int{2,1}
    fmt.Println("brute", brute(a,b))
}
