package main

import "fmt"

var (
   n    int
   adj  [][]int
   ansX []int64
   ansY []int64
   fx   = [4][2]int{{1, 0}, {0, 1}, {0, -1}, {-1, 0}}
)

func dfs(x int, y int64, fa int, lf int) {
   j := 0
   for _, v := range adj[x] {
       if v == fa {
           continue
       }
       if lf == j {
           j++
       }
       ansX[v] = ansX[x] + int64(fx[j][0])*y
       ansY[v] = ansY[x] + int64(fx[j][1])*y
       dfs(v, y>>1, x, 3-j)
       j++
   }
}

func main() {
   fmt.Scan(&n)
   adj = make([][]int, n+1)
   deg := make([]int, n+1)
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Scan(&u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
       deg[u]++
       deg[v]++
   }
   for i := 1; i <= n; i++ {
       if deg[i] > 4 {
           fmt.Println("NO")
           return
       }
   }
   ansX = make([]int64, n+1)
   ansY = make([]int64, n+1)
   ansX[1], ansY[1] = 0, 0
   dfs(1, 1<<30, 0, -1)
   fmt.Println("YES")
   for i := 1; i <= n; i++ {
       fmt.Println(ansX[i], ansY[i])
   }
}
