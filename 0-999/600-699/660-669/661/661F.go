package main
import(
    "bufio"
    "fmt"
    "os"
)
func sieve(n int) []bool{
    prime:=make([]bool,n+1)
    for i:=2;i<=n;i++{prime[i]=true}
    for p:=2;p*p<=n;p++{
        if prime[p]{
            for j:=p*p;j<=n;j+=p{prime[j]=false}
        }
    }
    return prime
}
func main(){
    primes:=sieve(500)
    in:=bufio.NewReader(os.Stdin)
    out:=bufio.NewWriter(os.Stdout)
    defer out.Flush()
    var t int
    if _,err:=fmt.Fscan(in,&t); err!=nil{ return }
    for i:=0;i<t;i++{
        var l,r int
        fmt.Fscan(in,&l,&r)
        cnt:=0
        for x:=l;x<=r;x++{ if primes[x]{cnt++} }
        fmt.Fprintln(out,cnt)
    }
}
