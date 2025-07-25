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
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, k int
   fmt.Fscan(reader, &n, &k)
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       grid[i] = []byte(s)
   }
   inf := n*n + 5
   dist := make([][]int, n)
   for i := 0; i < n; i++ {
       dist[i] = make([]int, n)
       for j := 0; j < n; j++ {
           dist[i][j] = inf
       }
   }
   // compute minimum non-'a' changes to each cell
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           cost := 0
           if grid[i][j] != 'a' {
               cost = 1
           }
           if i == 0 && j == 0 {
               dist[i][j] = cost
           } else {
               if i > 0 {
                   dist[i][j] = min(dist[i][j], dist[i-1][j]+cost)
               }
               if j > 0 {
                   dist[i][j] = min(dist[i][j], dist[i][j-1]+cost)
               }
           }
       }
   }
   // find farthest reachable prefix of all 'a's
   maxS := 0
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if dist[i][j] <= k {
               s := i + j + 2
               if s > maxS {
                   maxS = s
               }
           }
       }
   }
   totalLen := 2*n - 1
   ans := make([]byte, 0, totalLen)
   var frontier [][2]int
   vis := make([][]int, n)
   for i := range vis {
       vis[i] = make([]int, n)
   }
   mark := 1
   if maxS > 0 {
       // prefix of 'a's
       for i := 0; i < maxS-1; i++ {
           ans = append(ans, 'a')
       }
       // initial frontier
       for i := 0; i < n; i++ {
           for j := 0; j < n; j++ {
               if dist[i][j] <= k && i+j+2 == maxS {
                   frontier = append(frontier, [2]int{i, j})
                   vis[i][j] = mark
               }
           }
       }
   } else {
       // no changes possible, start from (0,0)
       ans = append(ans, grid[0][0])
       frontier = append(frontier, [2]int{0, 0})
       vis[0][0] = mark
       maxS = 2
   }
   // if entire path is covered
   if len(ans) == totalLen {
       writer.Write(ans)
       return
   }
   // extend answer one step at a time
   for len(ans) < totalLen {
       mark++
       best := byte('z' + 1)
       // find minimal next character
       for _, p := range frontier {
           i, j := p[0], p[1]
           if i+1 < n {
               c := grid[i+1][j]
               if c < best {
                   best = c
               }
           }
           if j+1 < n {
               c := grid[i][j+1]
               if c < best {
                   best = c
               }
           }
       }
       // build next frontier
       next := make([][2]int, 0)
       for _, p := range frontier {
           i, j := p[0], p[1]
           if i+1 < n && vis[i+1][j] != mark && grid[i+1][j] == best {
               vis[i+1][j] = mark
               next = append(next, [2]int{i + 1, j})
           }
           if j+1 < n && vis[i][j+1] != mark && grid[i][j+1] == best {
               vis[i][j+1] = mark
               next = append(next, [2]int{i, j + 1})
           }
       }
       ans = append(ans, best)
       frontier = next
   }
   writer.Write(ans)
}
