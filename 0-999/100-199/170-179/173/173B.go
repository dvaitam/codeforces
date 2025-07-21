package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       line := make([]byte, m)
       if _, err := fmt.Fscan(reader, &line); err != nil {
           return
       }
       grid[i] = line
   }
   const INF = 1e9
   // dp1H, dp1V: from person (1,1), directions: H=horizontal(east), V=vertical(down)
   dp1H := make([][]int, n)
   dp1V := make([][]int, n)
   for i := 0; i < n; i++ {
       dp1H[i] = make([]int, m)
       dp1V[i] = make([]int, m)
       for j := 0; j < m; j++ {
           dp1H[i][j] = INF
           dp1V[i][j] = INF
       }
   }
   dp1H[0][0] = 0
   // Fill dp1
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if i == 0 && j == 0 {
               continue
           }
           // horizontal arrival (from left)
           if j > 0 {
               // straight
               dp1H[i][j] = min(dp1H[i][j], dp1H[i][j-1])
               // turn at (i,j-1) from vertical to horizontal
               if grid[i][j-1] == '#' && dp1V[i][j-1] < INF {
                   dp1H[i][j] = min(dp1H[i][j], dp1V[i][j-1]+1)
               }
           }
           // vertical arrival (from above)
           if i > 0 {
               // straight
               dp1V[i][j] = min(dp1V[i][j], dp1V[i-1][j])
               // turn at (i-1,j) from horizontal to vertical
               if grid[i-1][j] == '#' && dp1H[i-1][j] < INF {
                   dp1V[i][j] = min(dp1V[i][j], dp1H[i-1][j]+1)
               }
           }
       }
   }
   ds := make([][]int, n)
   for i := 0; i < n; i++ {
       ds[i] = make([]int, m)
       for j := 0; j < m; j++ {
           ds[i][j] = min(dp1H[i][j], dp1V[i][j])
       }
   }
   // dp2H, dp2V: from basilisk (n,m), directions: H=horizontal(west), V=vertical(up)
   dp2H := make([][]int, n)
   dp2V := make([][]int, n)
   for i := 0; i < n; i++ {
       dp2H[i] = make([]int, m)
       dp2V[i] = make([]int, m)
       for j := 0; j < m; j++ {
           dp2H[i][j] = INF
           dp2V[i][j] = INF
       }
   }
   dp2H[n-1][m-1] = 0
   // Fill dp2
   for i := n - 1; i >= 0; i-- {
       for j := m - 1; j >= 0; j-- {
           if i == n-1 && j == m-1 {
               continue
           }
           // horizontal arrival (from right)
           if j < m-1 {
               // straight
               dp2H[i][j] = min(dp2H[i][j], dp2H[i][j+1])
               // turn at (i,j+1) from vertical to horizontal
               if grid[i][j+1] == '#' && dp2V[i][j+1] < INF {
                   dp2H[i][j] = min(dp2H[i][j], dp2V[i][j+1]+1)
               }
           }
           // vertical arrival (from below)
           if i < n-1 {
               // straight
               dp2V[i][j] = min(dp2V[i][j], dp2V[i+1][j])
               // turn at (i+1,j) from horizontal to vertical
               if grid[i+1][j] == '#' && dp2H[i+1][j] < INF {
                   dp2V[i][j] = min(dp2V[i][j], dp2H[i+1][j]+1)
               }
           }
       }
   }
   db := make([][]int, n)
   for i := 0; i < n; i++ {
       db[i] = make([]int, m)
       for j := 0; j < m; j++ {
           db[i][j] = min(dp2H[i][j], dp2V[i][j])
       }
   }
   ans := INF
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '#' && ds[i][j] < INF && db[i][j] < INF {
               cost := ds[i][j] + db[i][j]
               if cost < ans {
                   ans = cost
               }
           }
       }
   }
   if ans >= INF {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}
