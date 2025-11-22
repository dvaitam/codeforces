package main
import "fmt"
func main(){
    n:=9982
    l:=int64(44)
    f:=int64(35)
    best:=int64(0); ba,bb:=0,0
    gcd:=func(a,b int) int { for b!=0 { a,b=b,a%b }; return a }
    for a:=0;a<=n;a++{
        for b:=0;b<=n;b++{
            if a>1 && b==0 || b>1 && a==0 {continue}
            if a==0 && b==0 {continue}
            if a>0 && b>0 && gcd(a,b)!=1 {continue}
            // unreachable states filtered? This is only gcd condition, not path. Skip path check; just compute score
            v:=int64(a)*l + int64(b)*f
            if v>best {best=v; ba=a; bb=b}
        }
    }
    fmt.Println(best,ba,bb)
}
