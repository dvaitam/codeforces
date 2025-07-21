package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   grid := make([]string, 3)
   for i := 0; i < 3; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   // locate O
   var r0, c0 int
   for i := 0; i < 3; i++ {
       for j := 0; j < n; j++ {
           if grid[i][j] == 'O' {
               r0, c0 = i, j
           }
       }
   }
   // forbidden masks
   f := make([]int, n+1)
   for j := 0; j < n; j++ {
       m := 0
       for i := 0; i < 3; i++ {
           if grid[i][j] == 'X' || grid[i][j] == 'O' {
               m |= 1 << i
           }
       }
       f[j] = m
   }
   f[n] = 7
   // determine allowed moves by geometry
   allowL, allowR, allowV := false, false, false
   if c0 >= 2 {
       if grid[r0][c0-1] != 'X' && grid[r0][c0-2] != 'X' {
           allowL = true
       }
   }
   if c0+2 < n {
       if grid[r0][c0+1] != 'X' && grid[r0][c0+2] != 'X' {
           allowR = true
       }
   }
   if r0 == 0 {
       if grid[1][c0] != 'X' && grid[2][c0] != 'X' {
           allowV = true
       }
   } else if r0 == 2 {
       if grid[0][c0] != 'X' && grid[1][c0] != 'X' {
           allowV = true
       }
   }
   // DP arrays
   var dpPrev [8][8]int
   dpPrev[0][0] = 1
   // directions bits: L=0, R=1, V=2
   for i := 0; i < n; i++ {
       var dpNext [8][8]int
       // for each mask and mask2
       for mask := 0; mask < 8; mask++ {
           if mask & f[i] != 0 {
               continue
           }
           for m2 := 0; m2 < 8; m2++ {
               ways := dpPrev[mask][m2]
               if ways == 0 {
                   continue
               }
               // DFS fill column i
               var dfs func(pos, occ, nxt, m2cur int)
               dfs = func(pos, occ, nxt, m2cur int) {
                   if pos == 3 {
                       dpNext[nxt][m2cur] = (dpNext[nxt][m2cur] + ways) % mod
                       return
                   }
                   bit := 1 << pos
                   if occ & bit != 0 {
                       dfs(pos+1, occ, nxt, m2cur)
                       return
                   }
                   // try vertical
                   if pos+1 < 3 {
                       b2 := 1 << (pos + 1)
                       if occ & b2 == 0 {
                           m2v := m2cur
                           // check V
                           if allowV && i == c0 && ((r0 == 0 && pos == 1) || (r0 == 2 && pos == 0)) {
                               m2v |= 1 << 2
                           }
                           dfs(pos+2, occ|bit|b2, nxt, m2v)
                       }
                   }
                   // try horizontal
                   // to next column
                   if (f[i+1] & bit) == 0 && (nxt & bit) == 0 {
                       m2h := m2cur
                       if pos == r0 {
                           if allowL && i == c0-2 {
                               m2h |= 1 << 0
                           }
                           if allowR && i == c0+1 {
                               m2h |= 1 << 1
                           }
                       }
                       dfs(pos+1, occ|bit, nxt|bit, m2h)
                   }
               }
               dfs(0, mask|f[i], 0, m2)
           }
       }
       dpPrev = dpNext
   }
   // sum results
   var res int
   // mask must be 0
   Smask := 0
   if allowL {
       Smask |= 1 << 0
   }
   if allowR {
       Smask |= 1 << 1
   }
   if allowV {
       Smask |= 1 << 2
   }
   if Smask != 0 {
       for m2 := 0; m2 < 8; m2++ {
           if m2 & Smask != 0 {
               res = (res + dpPrev[0][m2]) % mod
           }
       }
   }
   fmt.Println(res)
}
