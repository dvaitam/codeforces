package main
import "fmt"
func gcd(a,b int) int { for b!=0 {a,b=b,a%b}; return a }
func main(){N:=30
 reach:=make([][]bool,N+1)
 for i:=range reach{reach[i]=make([]bool,N+1)}
 reach[1][0]=true; reach[0][1]=true
 for y:=1;y<=N;y++{reach[1][y]=true}
 for x:=1;x<=N;x++{reach[x][1]=true}
 for x:=2;x<=N;x++{
  for y:=2;y<=N;y++{
   if gcd(x,y)!=1 {continue}
   if gcd(x-1,y)==1 && reach[x-1][y]{reach[x][y]=true}
   if gcd(x,y-1)==1 && reach[x][y-1]{reach[x][y]=true}
  }
 }
 for x:=2;x<=N;x++{
  for y:=2;y<=N;y++{
   if gcd(x,y)!=1 || reach[x][y]{continue}
   left:=gcd(x-1,y)==1
   down:=gcd(x,y-1)==1
   if left||down{
    fmt.Println("unreachable with coprime neighbor",x,y,"left",left,"down",down)
    return
   }
  }
 }
 fmt.Println("none found upto",N)
}
