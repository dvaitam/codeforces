package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct {
   val   int64
   prevJ int
}

var (
   n       int
   graph   [][]bool
   matchY  []int
   matchX  []bool
)

func dfsMatch(u int, used []bool) bool {
   if used[u] {
       return false
   }
   used[u] = true
   for v := 0; v < n; v++ {
       if graph[u][v] {
           if matchY[v] < 0 || dfsMatch(matchY[v], used) {
               matchY[v] = u
               return true
           }
       }
   }
   return false
}

func dfsCover(u int, visL []bool, visR []bool) {
   if visL[u] {
       return
   }
   visL[u] = true
   for v := 0; v < n; v++ {
       if graph[u][v] && !visR[v] {
           visR[v] = true
           if matchY[v] >= 0 {
               dfsCover(matchY[v], visL, visR)
           }
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var m, l int
   fmt.Fscan(reader, &n, &m, &l)
   graph = make([][]bool, n)
   for i := 0; i < n; i++ {
       graph[i] = make([]bool, n)
   }
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       a--
       b--
       graph[a][b] = true
   }
   matchY = make([]int, n)
   for i := range matchY {
       matchY[i] = -1
   }
   matchX = make([]bool, n)
   // maximum matching
   for u := 0; u < n; u++ {
       used := make([]bool, n)
       if dfsMatch(u, used) {
           matchX[u] = true
       }
   }
   // vertex cover via König's theorem
   visL := make([]bool, n)
   visR := make([]bool, n)
   for u := 0; u < n; u++ {
       if !matchX[u] {
           dfsCover(u, visL, visR)
       }
   }
   // build cuts
   cuts := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if !visL[i] {
           cuts = append(cuts, i+1)
       }
       if visR[i] {
           cuts = append(cuts, -(i + 1))
       }
   }
   cLen := len(cuts)
   // dp[i][j]: i days used, j cuts remaining index
   dp := make([][]pair, l+1)
   for i := 0; i <= l; i++ {
       dp[i] = make([]pair, n+1)
       for j := 0; j <= n; j++ {
           dp[i][j].val = -1
           dp[i][j].prevJ = -1
       }
   }
   dp[0][cLen].val = 0
   // operations
   for i := 1; i <= l; i++ {
       var x, y int64
       fmt.Fscan(reader, &x, &y)
       // update dp
       for j := 0; j < n-i; j++ {
           for k := j; k <= n; k++ {
               prev := dp[i-1][k].val
               if prev >= 0 {
                   gain := x - int64(k-j)*y
                   if gain < 0 {
                       gain = 0
                   }
                   tot := prev + gain
                   if tot > dp[i][j].val {
                       dp[i][j].val = tot
                       dp[i][j].prevJ = k
                   }
               }
           }
       }
   }
   // find best
   lastJ := -1
   var best pair
   for j := 0; j <= n; j++ {
       if dp[l][j].val >= 0 && dp[l][j].val > best.val {
           best.val = dp[l][j].val
           lastJ = j
       }
   }
   // reconstruct
   seq := make([]int, 0)
   for i := l; i >= 1; i-- {
       prevJ := dp[i][lastJ].prevJ
       // day action
       seq = append(seq, 0)
       // cuts removed between lastJ and prevJ
       for lastJ < prevJ {
           seq = append(seq, cuts[lastJ])
           lastJ++
       }
   }
   // output reversed sequence
   total := len(seq)
   fmt.Fprintln(writer, total)
   for i := total - 1; i >= 0; i-- {
       if i > 0 {
           fmt.Fprint(writer, seq[i], " ")
       } else {
           fmt.Fprintln(writer, seq[i])
       }
   }
}
