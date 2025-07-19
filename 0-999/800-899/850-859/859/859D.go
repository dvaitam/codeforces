package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

var (
   n    int
   a    [][]float64
   v    [][][2]int
   dp   [][][2]float64
   used [][]bool
)

func getDp(x, b int) (float64, float64) {
   if used[x][b] {
       return dp[x][b][0], dp[x][b][1]
   }
   if x == 0 {
       // base: no rounds, score=0, probability=1
       used[x][b] = true
       dp[x][b][0] = 0
       dp[x][b][1] = 1
       return 0, 1
   }
   // get previous dp for same b
   pF, pS := getDp(x-1, b)
   var bestF, sumS float64
   l, r := v[x][b][0], v[x][b][1]
   for i := l; i <= r; i++ {
       // skip same segment in previous round
       if v[x-1][b][0] <= i && i <= v[x-1][b][1] {
           continue
       }
       tmpF, tmpS := getDp(x-1, i)
       sumS += tmpS * a[b][i]
       bestF = math.Max(bestF, tmpF+pF)
   }
   // adjust sumS by dividing by 100 (percent) and multiply by pS
   // since a[b][i] stored as probability (fraction)
   sumS *= pS
   // add score from this round: probability to advance * points
   bestF += sumS * float64(1<<(x-1))
   used[x][b] = true
   dp[x][b][0] = bestF
   dp[x][b][1] = sumS
   return bestF, sumS
}

func buildSegments(level, l, r int) {
   if level < 0 {
       return
   }
   for i := l; i <= r; i++ {
       v[level][i][0] = l
       v[level][i][1] = r
   }
   if level == 0 {
       return
   }
   mid := (l + r) / 2
   buildSegments(level-1, l, mid)
   buildSegments(level-1, mid+1, r)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read n
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   size := 1 << n
   // initialize arrays
   a = make([][]float64, size)
   for i := range a {
       a[i] = make([]float64, size)
       for j := 0; j < size; j++ {
           var p int
           fmt.Fscan(reader, &p)
           a[i][j] = float64(p) / 100.0
       }
   }
   v = make([][][2]int, n+1)
   dp = make([][][2]float64, n+1)
   used = make([][]bool, n+1)
   for lvl := 0; lvl <= n; lvl++ {
       v[lvl] = make([][2]int, size)
       dp[lvl] = make([][2]float64, size)
       used[lvl] = make([]bool, size)
   }
   buildSegments(n, 0, size-1)
   // compute answer
   var ans float64
   for i := 0; i < size; i++ {
       f, _ := getDp(n, i)
       if f > ans {
           ans = f
       }
   }
   // output with precision
   fmt.Printf("%.13f", ans)
}
