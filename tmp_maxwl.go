package main
import "fmt"
func gcd(a,b int) int { for b!=0 { a,b=b,a%b }; return a }
func main(){
    n:=27
    reach := make([][]bool,n+1)
    for i:=range reach{reach[i]=make([]bool,n+1)}
    reach[0][0]=true
    for a:=0;a<=n;a++{
        for b:=0;b<=n;b++{
            if !reach[a][b]{continue}
            if a<n && gcd(a+1,b)<=1 { reach[a+1][b]=true }
            if b<n && gcd(a,b+1)<=1 { reach[a][b+1]=true }
        }
    }
    maxA:=0
    for a:=0;a<=n;a++{
        if reach[a][27]{ maxA=a }
    }
    fmt.Println(maxA)
}
