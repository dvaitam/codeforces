package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // Read table, 1-based
   T := make([][]byte, n+1)
   for i := 1; i <= n; i++ {
       line := make([]byte, 0)
       var s string
       fmt.Fscan(in, &s)
       line = append(line, s...)
       // pad to 1-based
       T[i] = make([]byte, n+1)
       for j := 1; j <= n; j++ {
           T[i][j] = line[j-1]
       }
   }
   maxD := 2*n
   // positions on each diagonal d: r+c == d
   type pair struct{ r, c int }
   pos := make([][]pair, maxD+2)
   for d := 2; d <= maxD; d++ {
       for r := 1; r <= n; r++ {
           c := d - r
           if c >= 1 && c <= n {
               pos[d] = append(pos[d], pair{r, c})
           }
       }
   }
   // successors: for each diag d and index i in pos[d], list of indices in pos[d+1]
   succ := make([][][]int, maxD+2)
   for d := 2; d < maxD; d++ {
       m := len(pos[d])
       succ[d] = make([][]int, m)
       // map next positions
       nextMap := make(map[int]map[int]int)
       for j, p := range pos[d+1] {
           if nextMap[p.r] == nil {
               nextMap[p.r] = make(map[int]int)
           }
           nextMap[p.r][p.c] = j
       }
       for i, p := range pos[d] {
           // down: (r+1, c)
           if p.r+1 <= n {
               if idx, ok := nextMap[p.r+1][p.c]; ok {
                   succ[d][i] = append(succ[d][i], idx)
               }
           }
           // right: (r, c+1)
           if p.c+1 <= n {
               if idx, ok := nextMap[p.r][p.c+1]; ok {
                   succ[d][i] = append(succ[d][i], idx)
               }
           }
       }
   }
   // special case n==1: only one move
   if n == 1 {
       if T[1][1] == 'a' {
           fmt.Fprintln(out, "FIRST")
       } else if T[1][1] == 'b' {
           fmt.Fprintln(out, "SECOND")
       } else {
           fmt.Fprintln(out, "DRAW")
       }
       return
   }
   // dp[d][mask] = best score from diag d with current mask
   dp := make([]map[int]int, maxD+2)
   for d := 0; d <= maxD+1; d++ {
       dp[d] = make(map[int]int)
   }
   var solve func(d, mask int) int
   solve = func(d, mask int) int {
       // if beyond last diagonal, no more score
       if d > maxD-1 {
           return 0
       }
       if v, ok := dp[d][mask]; ok {
           return v
       }
       // determine current player: at diag d, move number is k=d-1
       turnFirst := (d%2) == 0
       var best int
       if turnFirst {
           best = -1e9
       } else {
           best = 1e9
       }
       // consider letters a-z for current positions on diag d
       for ch := byte('a'); ch <= 'z'; ch++ {
           // filter current positions by letter ch
           mask2 := 0
           for i := 0; i < len(pos[d]); i++ {
               if mask&(1<<i) != 0 {
                   p := pos[d][i]
                   if T[p.r][p.c] == ch {
                       mask2 |= 1 << i
                   }
               }
           }
           if mask2 == 0 {
               continue
           }
           // move to next diagonal
           nextMask := 0
           for i := 0; i < len(pos[d]); i++ {
               if mask2&(1<<i) != 0 {
                   for _, j := range succ[d][i] {
                       nextMask |= 1 << j
                   }
               }
           }
           var score int
           if ch == 'a' {
               score = 1
           } else if ch == 'b' {
               score = -1
           }
           val := score + solve(d+1, nextMask)
           if turnFirst {
               if val > best {
                   best = val
               }
           } else {
               if val < best {
                   best = val
               }
           }
       }
       dp[d][mask] = best
       return best
   }
   // initial move at (1,1) is forced
   startScore := 0
   if T[1][1] == 'a' {
       startScore = 1
   } else if T[1][1] == 'b' {
       startScore = -1
   }
   // initial mask for diag 3: successors of (1,1)
   initialMask := 0
   for _, j := range succ[2][0] {
       initialMask |= 1 << j
   }
   // solve from diag 3
   res := startScore + solve(3, initialMask)
   if res > 0 {
       fmt.Fprintln(out, "FIRST")
   } else if res < 0 {
       fmt.Fprintln(out, "SECOND")
   } else {
       fmt.Fprintln(out, "DRAW")
   }
}
