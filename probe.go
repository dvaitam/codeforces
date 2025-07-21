package main
import ("fmt")
func main(){k:=3;n:=20;dp:=make([][]bool,n+1);
 for i:=range dp{dp[i]=make([]bool,n+1)}
 for x:=0;x<=n;x++{
  for y:=0;y<=n;y++{
    if x==0 && y==0{dp[x][y]=false;continue}
    win:=false
    if x>0 && !dp[x-1][y]{win=true}
    if y>0 && !dp[x][y-1]{win=true}
    if !win && x>=k && y>=k && !dp[x-k][y-k]{win=true}
    dp[x][y]=win
  }
 }
 var P [][2]int
 for x:=0;x<=n;x++{ for y:=0;y<=n;y++{ if !dp[x][y]{P=append(P,[2]int{x,y})}}}
 fmt.Println("P-positions k=3 upto 20:")
 fmt.Println(P)
}
