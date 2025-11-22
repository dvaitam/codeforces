package main
import "fmt"
func gcd(a,b int) int { for b!=0 { a,b=b,a%b }; return a }
func brute(n,l,f int) int64 {
    reach := make([][]bool,n+1)
    for i:=range reach {reach[i]=make([]bool,n+1)}
    reach[0][0]=true
    best:=int64(0)
    for a:=0;a<=n;a++{
        for b:=0;b<=n;b++{
            if !reach[a][b]{continue}
            if (a<=1 && b==0)||(b<=1 && a==0)||(a>0 && b>0 && gcd(a,b)==1)||(a==0&&b==0){
                v:=int64(a)*int64(l)+int64(b)*int64(f)
                if v>best {best=v}
            }
            if a<n && gcd(a+1,b)<=1 {reach[a+1][b]=true}
            if b<n && gcd(a,b+1)<=1 {reach[a][b+1]=true}
        }
    }
    return best
}
func main(){
    for n:=2;n<=12;n++{
        // use l=f=1
        b:=brute(n,1,1)
        fmt.Println(n,b)
    }
}
