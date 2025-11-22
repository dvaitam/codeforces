package main
import "fmt"
func sieve(maxN int) []bool {
    is := make([]bool,maxN+1)
    for i:=2;i<=maxN;i++{is[i]=true}
    for i:=2;i*i<=maxN;i++{if is[i]{for j:=i*i;j<=maxN;j+=i{is[j]=false}}}
    return is
}
func main(){
    n:=9982
    is:=sieve(n)
    for i:=9970;i<=n;i++{
        if is[i]{fmt.Println(i)}
    }
}
