package main
import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var (

tWriter *bufio.Writer
)
func main(){
    in := bufio.NewReader(os.Stdin)
    t,T := readInt(in), readInt(in)
    n,m := readInt(in), readInt(in)
    li := make([]int64, n)
    ri := make([]int64, n)
    for i:=0;i<n;i++{
        li[i], ri[i] = readInt64(in), readInt64(in)
    }
    g := make([][]int, n)
    for i:=0;i<m;i++{
        u,v := readInt(in)-1, readInt(in)-1
        g[u] = append(g[u], v)
        g[v] = append(g[v], u)
    }
    color := make([]int, n)
    comp := make([]int, n)
    var comps [][]int
    var compCol0, compCol1 [][]int
    for i:=0;i<n;i++{
        if comp[i]!=0 {continue}
        // BFS
        q:=[]int{i}; comp[i]=len(comps)+1; color[i]=0
        col0, col1 := []int{}, []int{}
        col0 = append(col0, i)
        for k:=0;k<len(q);k++{
            u:=q[k]
            for _,v:=range g[u]{
                if comp[v]==0{
                    comp[v]=comp[u]
                    color[v]=color[u]^1
                    if color[v]==0 {col0 = append(col0, v)} else {col1 = append(col1, v)}
                    q = append(q, v)
                } else {
                    if comp[v]==comp[u] && color[v]==color[u] {
                        fmt.Println("IMPOSSIBLE")
                        return
                    }
                }
            }
        }
        comps = append(comps, append(col0,col1...)) // not used
        compCol0 = append(compCol0, col0)
        compCol1 = append(compCol1, col1)
    }
    C := len(compCol0)
    type CI struct{LA,RA,LB,RB int64; idx int}
    cis := make([]CI, C)
    for i:=0;i<C;i++{
        LA,RA := int64(0), int64(1e18)
        for _,u:=range compCol0[i]{
            if li[u]>LA {LA=li[u]}
            if ri[u]<RA {RA=ri[u]}
        }
        LB,RB := int64(0), int64(1e18)
        for _,u:=range compCol1[i]{
            if li[u]>LB {LB=li[u]}
            if ri[u]<RB {RB=ri[u]}
        }
        if LA>RA || LB>RB {
            fmt.Println("IMPOSSIBLE")
            return
        }
        cis[i]=CI{LA,RA,LB,RB,i}
    }
    // sort by LA-LB
    sort.Slice(cis, func(i,j int) bool{return cis[i].LA-cis[i].LB < cis[j].LA-cis[j].LB})
    // prefix/suffix
    prefLA := make([]int64, C+1); prefLB := make([]int64, C+1)
    prefRA := make([]int64, C+1); prefRB := make([]int64, C+1)
    const INF = int64(1e18)
    prefRA[0], prefRB[0] = INF, INF
    for i:=1;i<=C;i++{
        c := cis[i-1]
        prefLA[i] = max(prefLA[i-1], c.LA)
        prefLB[i] = max(prefLB[i-1], c.LB)
        prefRA[i] = min(prefRA[i-1], c.RA)
        prefRB[i] = min(prefRB[i-1], c.RB)
    }
    sufLA := make([]int64, C+2); sufLB := make([]int64, C+2)
    sufRA := make([]int64, C+2); sufRB := make([]int64, C+2)
    sufRA[C+1], sufRB[C+1] = INF, INF
    for i:=C;i>=1;i--{
        c := cis[i-1]
        sufLA[i] = max(sufLA[i+1], c.LA)
        sufLB[i] = max(sufLB[i+1], c.LB)
        sufRA[i] = min(sufRA[i+1], c.RA)
        sufRB[i] = min(sufRB[i+1], c.RB)
    }
    // try k
    var bestK int = -1
    var outL1,outL2 int64
    for k:=0;k<=C;k++{
        L1 := max(prefLA[k], sufLB[k+1])
        R1 := min(prefRA[k], sufRB[k+1])
        L2 := max(prefLB[k], sufLA[k+1])
        R2 := min(prefRB[k], sufRA[k+1])
        if L1>R1 || L2>R2 {continue}
        if L1+L2 > T || R1+R2 < t {continue}
        // feasible
        bestK = k; outL1=L1; outL2=L2; break
    }
    if bestK<0 {
        fmt.Println("IMPOSSIBLE")
        return
    }
    // choose n1,n2
    n1 := outL1; n2 := outL2
    if n1+n2 < t {
        delta := t - (n1+n2)
        add := min(delta, (min(prefRA[bestK], sufRB[bestK+1]) - n1))
        n1 += add; delta -= add
        add = min(delta, (min(prefRB[bestK], sufRA[bestK+1]) - n2))
        n2 += add; delta -= add
    }
    // print
    fmt.Println("POSSIBLE")
    fmt.Printf("%d %d\n", n1, n2)
    // reconstruct assignment
    ans := make([]byte, n)
    for i:=0;i<n;i++{ ans[i]='1' }
    for idx:=0;idx<bestK;idx++{
        ci := cis[idx]
        for _,u:=range compCol0[ci.idx] { ans[u]='1' }
        for _,u:=range compCol1[ci.idx] { ans[u]='2' }
    }
    for idx:=bestK;idx<C;idx++{
        ci := cis[idx]
        for _,u:=range compCol0[ci.idx] { ans[u]='2' }
        for _,u:=range compCol1[ci.idx] { ans[u]='1' }
    }
    fmt.Println(string(ans))
}

func readInt(r *bufio.Reader) int {
    var x int
    fmt.Fscan(r, &x)
    return x
}
func readInt64(r *bufio.Reader) int64 {
    var x int64
    fmt.Fscan(r, &x)
    return x
}
func max(a,b int64) int64 { if a>b {return a} ; return b }
func min(a,b int64) int64 { if a<b {return a} ; return b }
