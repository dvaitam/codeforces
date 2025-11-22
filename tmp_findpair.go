package main
import "fmt"
func gcd(a,b int) int { for b!=0 { a,b=b,a%b }; return a }
func main(){
    n:=27
    l:=int64(4)
    f:=int64(5)
    reach := make([][]bool,n+1)
    for i:=range reach{reach[i]=make([]bool,n+1)}
    reach[0][0]=true
    best:=int64(0); ba,bb:=0,0
    for a:=0;a<=n;a++{
        for b:=0;b<=n;b++{
            if !reach[a][b]{continue}
            if (a<=1 && b==0)||(b<=1 && a==0)||(a>0 && b>0 && gcd(a,b)==1)||(a==0&&b==0){
                v:=int64(a)*l + int64(b)*f
                if v>best { best=v; ba=a; bb=b }
            }
            if a<n && gcd(a+1,b)<=1 { reach[a+1][b]=true }
            if b<n && gcd(a,b+1)<=1 { reach[a][b+1]=true }
        }
    }
    fmt.Println(best, ba, bb)
}
