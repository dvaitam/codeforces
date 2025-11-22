package main
import "fmt"
type node struct{a,b int}
func gcd(a,b int) int { for b!=0 { a,b=b,a%b }; return a }
func main(){
    n:=27
    targetA,targetB := 20,27
    reach := make([][]bool,n+1)
    pre := make([][]node,n+1)
    for i:=range reach{reach[i]=make([]bool,n+1); pre[i]=make([]node,n+1)}
    reach[0][0]=true
    queue := []node{{0,0}}
    dirs := []node{{1,0},{0,1}}
    for len(queue)>0{
        cur:=queue[0]; queue=queue[1:]
        for _,d := range dirs {
            na,nb := cur.a+d.a, cur.b+d.b
            if na>n || nb>n || reach[na][nb]{continue}
            if gcd(na,nb)<=1 {
                reach[na][nb]=true
                pre[na][nb]=cur
                queue=append(queue,node{na,nb})
            }
        }
    }
    if !reach[targetA][targetB]{fmt.Println("unreachable");return}
    path := []node{}
    a,b := targetA,targetB
    for !(a==0 && b==0){
        path = append(path,node{a,b})
        p := pre[a][b]
        a,b = p.a,p.b
    }
    path = append(path,node{0,0})
    for i:=len(path)-1;i>=0;i--{
        fmt.Printf("(%d,%d)\n", path[i].a, path[i].b)
    }
}
