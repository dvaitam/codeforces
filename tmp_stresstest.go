package main
import (
    "fmt"
    "math/rand"
    "time"
)
func sieve(maxN int) []int {
    spf := make([]int,maxN+2)
    primes := []int{}
    for i:=2;i<=maxN;i++{
        if spf[i]==0 { spf[i]=i; primes = append(primes, i) }
        for _,p := range primes {
            if p>spf[i] || i*p>maxN { break }
            spf[i*p]=p
        }
    }
    spf[1]=maxN+1
    return spf
}
func calc(n int, l,f int64, spf []int) int64 {
    maxSpf := make([]int, n+2)
    maxSpf[n+1]=1
    for i:=2;i<=n;i++{
        p:=spf[i]
        if i>maxSpf[p]{maxSpf[p]=i}
    }
    rough := make([]int, n+2)
    for q:=n; q>=1; q--{
        rough[q]=rough[q+1]
        if v:=maxSpf[q+1]; v>rough[q]{rough[q]=v}
    }
    ans := l*int64(n)+f
    if v:=f*int64(n)+l; v>ans { ans=v }
    for q:=2; q<=n; q++{
        r:=rough[q]
        block:=spf[q]
        limit:=r + (block - r%block -1)
        if limit>n { limit=n }
        if v:=l*int64(limit)+f*int64(q); v>ans { ans=v }
        if v:=f*int64(limit)+l*int64(q); v>ans { ans=v }
    }
    return ans
}
func brute(n int, l,f int64) int64 {
    reach := make([][]bool,n+1)
    for i:=range reach{reach[i]=make([]bool,n+1)}
    reach[0][0]=true
    gcd:=func(a,b int) int { for b!=0 { a,b=b,a%b }; return a }
    best:=int64(0)
    for a:=0;a<=n;a++{
        for b:=0;b<=n;b++{
            if !reach[a][b]{continue}
            if (a<=1 && b==0)||(b<=1 && a==0)||(a>0 && b>0 && gcd(a,b)==1)||(a==0&&b==0){
                v:=int64(a)*l + int64(b)*f
                if v>best { best=v }
            }
            if a<n && gcd(a+1,b)<=1 { reach[a+1][b]=true }
            if b<n && gcd(a,b+1)<=1 { reach[a][b+1]=true }
        }
    }
    return best
}
func main(){
    rand.Seed(time.Now().UnixNano())
    for n:=2;n<=30;n++{
        spf:=sieve(n)
        for t:=0;t<200;t++{
            l:=int64(rand.Intn(5)+1)
            f:=int64(rand.Intn(5)+1)
            b:=brute(n,l,f)
            c:=calc(n,l,f,spf)
            if b!=c {
                fmt.Printf("Mismatch n=%d l=%d f=%d b=%d c=%d\n", n,l,f,b,c)
                return
            }
        }
    }
    fmt.Println("all good up to 30 random tests")
}
