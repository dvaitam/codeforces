package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, d int
   if _, err := fmt.Fscan(reader, &n, &m, &d); err != nil {
       return
   }
   owner := make([]int, n+1)
   for i := 0; i < m; i++ {
       var sz int
       fmt.Fscan(reader, &sz)
       for j := 0; j < sz; j++ {
           var v int
           fmt.Fscan(reader, &v)
           owner[v] = i
       }
   }
   // number of windows
   nw := n - d + 1
   masks := make([]uint32, nw)
   cnt := make([]int, m)
   var curMask uint32
   // initial window [1..d]
   for i := 1; i <= d; i++ {
       oi := owner[i]
       if cnt[oi] == 0 {
           curMask |= 1 << oi
       }
       cnt[oi]++
   }
   masks[0] = curMask
   for i := 1; i < nw; i++ {
       // remove i, add i+d
       oi := owner[i]
       cnt[oi]--
       if cnt[oi] == 0 {
           curMask &^= 1 << oi
       }
       ni := owner[i+d]
       if cnt[ni] == 0 {
           curMask |= 1 << ni
       }
       cnt[ni]++
       masks[i] = curMask
   }
   // check if chosen mask covers all windows
   var validMask = func(chosen uint32) bool {
       for i := 0; i < nw; i++ {
           if masks[i]&chosen == 0 {
               return false
           }
       }
       return true
   }
   // try combinations by increasing k
   // recursive generate combinations
   var found bool
   var answer int
   var dfs func(start, depth, k int, chosen uint32)
   dfs = func(start, depth, k int, chosen uint32) {
       if found {
           return
       }
       if depth == k {
           if validMask(chosen) {
               found = true
               answer = k
           }
           return
       }
       // prune: need at least k-depth elements from start..m-1
       for i := start; i <= m-(k-depth); i++ {
           dfs(i+1, depth+1, k, chosen|(1<<uint(i)))
           if found {
               return
           }
       }
   }
   for k := 1; k <= m; k++ {
       dfs(0, 0, k, 0)
       if found {
           fmt.Println(answer)
           return
       }
   }
   // fallback: all sets
   fmt.Println(m)
}
