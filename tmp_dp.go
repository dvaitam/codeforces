package main
import "fmt"
const INF=int(1e18)
func solve(a,b []int) int {
    n:=len(a)
    R:=1
    changed:=true
    for changed {
        changed=false
        cur:=R
        for i:=1;i<=R;i++{
            if b[i-1]>cur {cur=b[i-1]; changed=true}
        }
        R=cur
    }
    pref:=make([]int,n+1)
    for i,v:=range a {pref[i+1]=pref[i]+v}
    if R==1 {return a[0]}
    size:=R
    seg:=make([]int,4*size)
    for i:=range seg {seg[i]=INF}
    var upd func(int,int,int,int,int,int)
    upd=func(node,l,r,ql,qr,val int){
        if ql<=l && r<=qr {
            if val<seg[node] {seg[node]=val}
            return
        }
        m:=(l+r)>>1
        if ql<=m {upd(node<<1,l,m,ql,qr,val)}
        if qr>m {upd(node<<1|1,m+1,r,ql,qr,val)}
    }
    var query func(int,int,int,int) int
    query=func(node,l,r,idx int) int{
        res:=seg[node]
        if l==r {return res}
        m:=(l+r)>>1
        var t int
        if idx<=m {t=query(node<<1,l,m,idx)} else {t=query(node<<1|1,m+1,r,idx)}
        if t<res {res=t}
        return res
    }
    dp:=make([]int,size)
    for i:=range dp {dp[i]=INF}
    dp[0]=0
    for i:=1;i<=R;i++{
        val:=dp[i-1]+a[i-1]
        l:=i
        r:=b[i-1]-1
        if r>R-1 {r=R-1}
        if l<=r {upd(1,1,R-1,l,r,val)}
        if i<=R-1 {dp[i]=query(1,1,R-1,i)}
    }
    minCost:=dp[R-1]
    return pref[R]-minCost
}
func main(){
    a:=[]int{15,16}
    b:=[]int{2,1}
    fmt.Println(solve(a,b))
}
