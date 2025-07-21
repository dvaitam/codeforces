package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   bannedRow := make([]bool, n+1)
   bannedCol := make([]bool, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       if x >= 1 && x <= n {
           bannedRow[x] = true
       }
       if y >= 1 && y <= n {
           bannedCol[y] = true
       }
   }
   // Determine clear rows and columns (excluding corners irrelevant)
   rowClear := make([]bool, n+1)
   colClear := make([]bool, n+1)
   for i := 2; i <= n-1; i++ {
       rowClear[i] = !bannedRow[i]
       colClear[i] = !bannedCol[i]
   }
   total := 0
   // Process components by pairing c and p=n+1-c
   // indices c from 2..n-1
   seen := make([]bool, n+1)
   for c := 2; c <= n-1; c++ {
       if seen[c] {
           continue
       }
       p := n + 1 - c
       seen[c] = true
       if p >= 2 && p <= n-1 {
           seen[p] = true
       }
       // collect nodes in this component
       type node struct { kind byte; idx int }
       var nodes []node
       // column c
       if colClear[c] {
           nodes = append(nodes, node{'T', c}, node{'B', c})
       }
       // row c
       if rowClear[c] {
           nodes = append(nodes, node{'L', c}, node{'R', c})
       }
       // if p!=c, add for p
       if p != c && p >= 2 && p <= n-1 {
           if colClear[p] {
               nodes = append(nodes, node{'T', p}, node{'B', p})
           }
           if rowClear[p] {
               nodes = append(nodes, node{'L', p}, node{'R', p})
           }
       }
       K := len(nodes)
       if K == 0 {
           continue
       }
       // build conflict matrix
       conflict := make([][]bool, K)
       for i := range conflict {
           conflict[i] = make([]bool, K)
       }
       for i := 0; i < K; i++ {
           for j := i + 1; j < K; j++ {
               a := nodes[i]
               b := nodes[j]
               c1, c2 := a.idx, b.idx
               conflictFlag := false
               // same column T/B
               if (a.kind == 'T' || a.kind == 'B') && (b.kind == 'T' || b.kind == 'B') && c1 == c2 {
                   conflictFlag = true
               }
               // same row L/R
               if (a.kind == 'L' || a.kind == 'R') && (b.kind == 'L' || b.kind == 'R') && c1 == c2 {
                   conflictFlag = true
               }
               // vertical a, horizontal b
               if (a.kind == 'T' || a.kind == 'B') && (b.kind == 'L' || b.kind == 'R') {
                   // top-left
                   if a.kind == 'T' && b.kind == 'L' && c1 == c2 {
                       conflictFlag = true
                   }
                   // bottom-right
                   if a.kind == 'B' && b.kind == 'R' && c1 == c2 {
                       conflictFlag = true
                   }
                   // top-right
                   if a.kind == 'T' && b.kind == 'R' && c1 + c2 == n+1 {
                       conflictFlag = true
                   }
                   // bottom-left
                   if a.kind == 'B' && b.kind == 'L' && c1 + c2 == n+1 {
                       conflictFlag = true
                   }
               }
               // horizontal a, vertical b (symmetric)
               if !conflictFlag && (a.kind == 'L' || a.kind == 'R') && (b.kind == 'T' || b.kind == 'B') {
                   // mirror above by swapping
                   if b.kind == 'T' && a.kind == 'L' && c1 == c2 {
                       conflictFlag = true
                   }
                   if b.kind == 'B' && a.kind == 'R' && c1 == c2 {
                       conflictFlag = true
                   }
                   if b.kind == 'T' && a.kind == 'R' && c1 + c2 == n+1 {
                       conflictFlag = true
                   }
                   if b.kind == 'B' && a.kind == 'L' && c1 + c2 == n+1 {
                       conflictFlag = true
                   }
               }
               conflict[i][j] = conflictFlag
               conflict[j][i] = conflictFlag
           }
       }
       // brute MIS
       best := 0
       var maskBest int
       for mask := 0; mask < (1 << K); mask++ {
           cnt := bitsOn(mask)
           if cnt <= best {
               continue
           }
           ok := true
           for i := 0; i < K && ok; i++ {
               if mask&(1<<i) == 0 {
                   continue
               }
               for j := i + 1; j < K; j++ {
                   if mask&(1<<j) != 0 && conflict[i][j] {
                       ok = false
                       break
                   }
               }
           }
           if ok {
               best = cnt
               maskBest = mask
           }
       }
       total += best
   }
   fmt.Println(total)
}

// bitsOn counts set bits
func bitsOn(x int) int {
   cnt := 0
   for x > 0 {
       cnt += x & 1
       x >>= 1
   }
   return cnt
}
