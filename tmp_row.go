package main
import "fmt"
func gcd(a,b int) int { for b!=0 { a,b=b,a%b }; return a }
func main(){
    n:=27; b:=27
    reach := make([][]bool,n+1)
    for i:=range reach{reach[i]=make([]bool,n+1)}
    reach[0][0]=true
    for a:=0;a<=n;a++{
        for bb:=0;bb<=n;bb++{
            if !reach[a][bb]{continue}
            if a<n && gcd(a+1,bb)<=1 { reach[a+1][bb]=true }
            if bb<n && gcd(a,bb+1)<=1 { reach[a][bb+1]=true }
        }
    }
    for a:=0;a<=n;a++{
        if reach[a][b] { fmt.Print(a," ") }
    }
    fmt.Println()
}
