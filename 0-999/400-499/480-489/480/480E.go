package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b, c int) int {
   if a < b {
       if a < c {
           return a
       }
       return c
   }
   if b < c {
       return b
   }
   return c
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)
   grid := make([][]bool, n)
   for i := 0; i < n; i++ {
       line, _ := reader.ReadString('\n')
       if len(line) < m {
           // maybe no newline
           buf := make([]byte, m+1)
           copy(buf, line)
           l, _ := reader.Read(buf[len(line):])
           line = string(buf[:len(line)+l])
       }
       row := make([]bool, m)
       for j := 0; j < m; j++ {
           row[j] = (line[j] == '.')
       }
       grid[i] = row
   }
   xs := make([]int, k)
   ys := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i])
       xs[i]--
       ys[i]--
       grid[xs[i]][ys[i]] = false
   }
   // dp of largest square ending at (i,j)
   dp := make([][]int, n)
   for i := range dp {
       dp[i] = make([]int, m)
   }
   maxSq := 0
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] {
               if i > 0 && j > 0 {
                   dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
               } else {
                   dp[i][j] = 1
               }
               if dp[i][j] > maxSq {
                   maxSq = dp[i][j]
               }
           }
       }
   }
   // answers at reverse times
   ans := make([]int, k+1)
   ans[k] = maxSq
   // queue for BFS
   type pt struct{ i, j int }
   q := make([]pt, 0, n*m)
   for t := k - 1; t >= 0; t-- {
       i := xs[t]
       j := ys[t]
       grid[i][j] = true
       // compute dp at (i,j)
       old := dp[i][j]
       var newv int
       if i > 0 && j > 0 {
           newv = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
       } else {
           newv = 1
       }
       if newv > old {
           dp[i][j] = newv
           if newv > maxSq {
               maxSq = newv
           }
           q = q[:0]
           q = append(q, pt{i, j})
           // BFS propagate
           for head := 0; head < len(q); head++ {
               p := q[head]
               // neighbors that depend on (p.i, p.j)
               for _, d := range [][2]int{{1, 0}, {0, 1}, {1, 1}} {
                   ni := p.i + d[0]
                   nj := p.j + d[1]
                   if ni < n && nj < m && grid[ni][nj] {
                       // recompute dp
                       var cand int
                       if ni > 0 && nj > 0 {
                           cand = 1 + min(dp[ni-1][nj], dp[ni][nj-1], dp[ni-1][nj-1])
                       } else {
                           cand = 1
                       }
                       if cand > dp[ni][nj] {
                           dp[ni][nj] = cand
                           if cand > maxSq {
                               maxSq = cand
                           }
                           q = append(q, pt{ni, nj})
                       }
                   }
               }
           }
       }
       ans[t] = maxSq
   }
   // output answers after each car enters: t=1..k
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i := 1; i <= k; i++ {
       fmt.Fprintln(w, ans[i])
   }
}
