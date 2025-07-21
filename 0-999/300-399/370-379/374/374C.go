package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m   int
   grid    [][]byte
   dp      [][]int
   vis     [][]int // 0=unvisited,1=visiting,2=visited
   infinite bool
   dirs    = [4][2]int{{-1,0},{1,0},{0,-1},{0,1}}
   nextCh  = map[byte]byte{'D':'I','I':'M','M':'A','A':'D'}
)

func dfs(i, j int) int {
   if infinite {
       return 0
   }
   if vis[i][j] == 1 {
       infinite = true
       return 0
   }
   if vis[i][j] == 2 {
       return dp[i][j]
   }
   vis[i][j] = 1
   best := 1
   cur := grid[i][j]
   nxt, ok := nextCh[cur]
   if ok {
       for _, d := range dirs {
           ni, nj := i+d[0], j+d[1]
           if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == nxt {
               v := dfs(ni, nj)
               if infinite {
                   return 0
               }
               if v+1 > best {
                   best = v + 1
               }
           }
       }
   }
   dp[i][j] = best
   vis[i][j] = 2
   return best
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m)
   grid = make([][]byte, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       grid[i] = []byte(line)
   }
   dp = make([][]int, n)
   vis = make([][]int, n)
   for i := 0; i < n; i++ {
       dp[i] = make([]int, m)
       vis[i] = make([]int, m)
   }
   maxLen := 0
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == 'D' {
               v := dfs(i, j)
               if infinite {
                   fmt.Fprintln(writer, "Poor Inna!")
                   return
               }
               if v > maxLen {
                   maxLen = v
               }
           }
       }
   }
   if maxLen < 4 {
       fmt.Fprintln(writer, "Poor Dima!")
   } else {
       fmt.Fprintln(writer, maxLen/4)
   }
}
