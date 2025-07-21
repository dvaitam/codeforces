package main

import (
   "bufio"
   "fmt"
   "os"
)

// Query holds a single query
type Query struct {
   x1, y1, x2, y2 int
   idx            int
   Pbits          []uint64
   ans            bool
}

var (
   n, m int
   grid [][]byte
   // rowWalls[i][j]: number of '#' in row i, cols [1..j]
   rowWalls [][]int
   // answers output
   answers []bool
   mWords   int
)

// solve processes queries for rows [l..r]
func solve(l, r int, qs []*Query) {
   if len(qs) == 0 {
       return
   }
   if l == r {
       // single row: only right moves
       for _, q := range qs {
           // check continuous white from y1 to y2
           if rowWalls[l][q.y2] - rowWalls[l][q.y1-1] == 0 {
               answers[q.idx] = true
           } else {
               answers[q.idx] = false
           }
       }
       return
   }
   mid := (l + r) >> 1
   // buckets
   startBuckets := make([][]*Query, mid-l+1)
   endBuckets := make([][]*Query, r-mid+1)
   var leftQ, rightQ []*Query
   for _, q := range qs {
       if q.x2 < mid {
           leftQ = append(leftQ, q)
       } else if q.x1 > mid {
           rightQ = append(rightQ, q)
       } else {
           // spans mid
           startBuckets[q.x1-l] = append(startBuckets[q.x1-l], q)
           endBuckets[q.x2-mid] = append(endBuckets[q.x2-mid], q)
       }
   }
   // top DP: from mid down to l
   // dpNext[j] is bitset of reachable mid columns from current row cell (i,j)
   // dpNext[j] and dpRow[j] are bitsets for j=0..m+1
   dpNext := make([][]uint64, m+2)
   dpRow := make([][]uint64, m+2)
   for j := 0; j <= m+1; j++ {
       dpNext[j] = make([]uint64, mWords)
       dpRow[j] = make([]uint64, mWords)
   }
   // init mid row in dpNext
   for j := 1; j <= m; j++ {
       if grid[mid][j] == '.' {
           bit := uint(j - 1)
           dpNext[j][bit/64] |= 1 << (bit % 64)
       }
   }
   // process rows
   for i := mid; i >= l; i-- {
       if i < mid {
           // compute dpRow for row i
           for j := m; j >= 1; j-- {
               if grid[i][j] == '.' {
                   // OR down and right
                   for d := 0; d < mWords; d++ {
                       dpRow[j][d] = dpNext[j][d] | dpRow[j+1][d]
                   }
               } else {
                   for d := 0; d < mWords; d++ {
                       dpRow[j][d] = 0
                   }
               }
           }
           // swap dpRow into dpNext
           for j := 1; j <= m; j++ {
               copy(dpNext[j], dpRow[j])
           }
       }
       // answer start queries at row i
       for _, q := range startBuckets[i-l] {
           // store Pbits for start
           q.Pbits = make([]uint64, mWords)
           copy(q.Pbits, dpNext[q.y1])
       }
   }
   // bottom DP: from mid up to r
   // reset dpNext and dpRow fully for bottom DP
   for j := 0; j <= m+1; j++ {
       for d := 0; d < mWords; d++ {
           dpNext[j][d] = 0
           dpRow[j][d] = 0
       }
   }
   // init mid row in dpNext
   for j := 1; j <= m; j++ {
       if grid[mid][j] == '.' {
           bit := uint(j - 1)
           dpNext[j][bit/64] |= 1 << (bit % 64)
       }
   }
   for i := mid; i <= r; i++ {
       if i > mid {
           for j := 1; j <= m; j++ {
               if grid[i][j] == '.' {
                   for d := 0; d < mWords; d++ {
                       dpRow[j][d] = dpNext[j][d] | dpRow[j-1][d]
                   }
               } else {
                   for d := 0; d < mWords; d++ {
                       dpRow[j][d] = 0
                   }
               }
           }
           for j := 1; j <= m; j++ {
               copy(dpNext[j], dpRow[j])
           }
       }
       // answer end queries at row i
       for _, q := range endBuckets[i-mid] {
           // intersect q.Pbits and dpNext[q.y2]
           lbit := q.y1 - 1
           rbit := q.y2 - 1
           w1 := lbit / 64
           w2 := rbit / 64
           ok := false
           for w := w1; w <= w2; w++ {
               mask := ^uint64(0)
               if w == w1 {
                   mask &= ^((1 << (lbit % 64)) - 1)
               }
               if w == w2 {
                   mask &= (1 << ((rbit % 64) + 1)) - 1
               }
               if q.Pbits[w]&dpNext[q.y2][w]&mask != 0 {
                   ok = true
                   break
               }
           }
           answers[q.idx] = ok
       }
   }
   // recurse
   solve(l, mid-1, leftQ)
   solve(mid+1, r, rightQ)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &m)
   grid = make([][]byte, n+1)
   var line string
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &line)
       grid[i] = []byte(" " + line)
   }
   // preprocess rowWalls
   rowWalls = make([][]int, n+1)
   for i := 1; i <= n; i++ {
       rowWalls[i] = make([]int, m+1)
       for j := 1; j <= m; j++ {
           rowWalls[i][j] = rowWalls[i][j-1]
           if grid[i][j] == '#' {
               rowWalls[i][j]++
           }
       }
   }
   var q int
   fmt.Fscan(in, &q)
   queries := make([]*Query, q)
   answers = make([]bool, q)
   // bitset words
   mWords = (m + 63) / 64
   for i := 0; i < q; i++ {
       var x1, y1, x2, y2 int
       fmt.Fscan(in, &x1, &y1, &x2, &y2)
       queries[i] = &Query{x1: x1, y1: y1, x2: x2, y2: y2, idx: i}
   }
   solve(1, n, queries)
   for i := 0; i < q; i++ {
       if answers[i] {
           fmt.Fprintln(out, "Yes")
       } else {
           fmt.Fprintln(out, "No")
       }
   }
}
