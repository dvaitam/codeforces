package main

import (
   "bufio"
   "fmt"
   "os"
)

var n, m, np int
var adj []int
var MEMadj []int
var MEM [][]int
var lnxt []int

// dp returns 1 if a valid routing path exists starting from u with bitmask bm
func dp(u, bm int) int {
   if MEM[u][bm] != -1 {
       return MEM[u][bm]
   }
   res := 0
   for i := 0; i < n; i++ {
       v := 1 << i
       if (v&bm) != 0 && (v&adj[u]) != 0 {
           if (v^bm) == 0 {
               res = 1
               lnxt[u] = i
               MEM[u][bm] = res
               return res
           }
           if bm%v != 0 {
               p := dp(i, bm^v)
               if p == 1 {
                   lnxt[u] = i
                   MEM[u][bm] = p
                   return p
               }
           }
       }
   }
   MEM[u][bm] = res
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m)
   np = 1 << n
   adj = make([]int, n)
   MEMadj = make([]int, np)
   MEM = make([][]int, n)
   for i := 0; i < n; i++ {
       MEM[i] = make([]int, np)
       for j := 0; j < np; j++ {
           MEM[i][j] = -1
       }
   }
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       a--
       b--
       adj[a] |= 1 << b
       adj[b] |= 1 << a
   }
   // Precompute neighbor mask for every subset
   for mask := 0; mask < np; mask++ {
       for j := 0; j < n; j++ {
           if (mask & (1 << j)) != 0 {
               MEMadj[mask] |= adj[j]
           }
       }
   }
   lnxt = make([]int, n)
   for i := 0; i < n; i++ {
       lnxt[i] = -1
   }
   ispos := false
   // Find a dominating subset and attempt to build routing
   for mask := 0; mask < np && !ispos; mask++ {
       if MEMadj[mask] == np-1 {
           for j := 0; j < n; j++ {
               if (mask&(1<<j)) != 0 {
                   if dp(j, mask) == 1 {
                       if lnxt[j] == j {
                           for k := 0; k < n; k++ {
                               if k == j {
                                   continue
                               }
                               if (adj[j] & (1 << k)) != 0 {
                                   lnxt[j] = k
                               }
                               break
                           }
                       }
                       ispos = true
                       break
                   }
                   break
               }
           }
       }
   }
   if !ispos {
       fmt.Fprintln(writer, "NO")
       return
   }
   // Assign remaining auxiliary servers
   for i := 0; i < n; i++ {
       if lnxt[i] == -1 {
           for j := 0; j < n; j++ {
               if (adj[i]&(1<<j)) != 0 && lnxt[j] != -1 {
                   lnxt[i] = j
                   break
               }
           }
       }
   }
   fmt.Fprintln(writer, "YES")
   for i := 0; i < n; i++ {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, lnxt[i]+1)
   }
   fmt.Fprintln(writer)
}
