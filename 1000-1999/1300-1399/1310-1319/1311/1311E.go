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
   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var n, d int
       fmt.Fscan(in, &n, &d)
       // max sum depths = n*(n-1)/2, min sum for balanced binary tree
       maxSum := n*(n-1)/2
       // compute minSum
       rem := n - 1
       last := 1
       depth := 1
       minSum := 0
       for rem > 0 {
           cnt := last * 2
           if cnt > rem {
               cnt = rem
           }
           minSum += cnt * depth
           rem -= cnt
           last = cnt
           depth++
       }
       if d < minSum || d > maxSum {
           fmt.Fprintln(out, "NO")
           continue
       }
       // initialize chain
       parent := make([]int, n+1)
       depthArr := make([]int, n+1)
       childCnt := make([]int, n+1)
       for i := 2; i <= n; i++ {
           parent[i] = i - 1
           childCnt[i-1]++
           depthArr[i] = i - 1
       }
       depthArr[1] = 0
       // available parents by depth
       avail := make([][]int, n)
       for i := 1; i <= n; i++ {
           cap := 2 - childCnt[i]
           if cap > 0 {
               d0 := depthArr[i]
               avail[d0] = append(avail[d0], i)
           }
       }
       // leaves by depth
       leaves := make([][]int, n)
       for i := 1; i <= n; i++ {
           if childCnt[i] == 0 {
               d0 := depthArr[i]
               leaves[d0] = append(leaves[d0], i)
           }
       }
       diff := maxSum - d
       // pointers for scanning
       leafMax := n - 1
       for diff > 0 {
           // find deepest leaf
           for leafMax > 0 && len(leaves[leafMax]) == 0 {
               leafMax--
           }
           u := leaves[leafMax][len(leaves[leafMax])-1]
           leaves[leafMax] = leaves[leafMax][:len(leaves[leafMax])-1]
           du := depthArr[u]
           // compute minimal parent depth
           dvMin := du - diff - 1
           if dvMin < 0 {
               dvMin = 0
           }
           // find parent depth
           var dv int
           var v int
           for dv = dvMin; dv < du; dv++ {
               if dv < len(avail) && len(avail[dv]) > 0 {
                   v = avail[dv][0]
                   break
               }
           }
           // new depth and reduction
           nd := dv + 1
           delta := du - nd
           // update diff
           diff -= delta
           // reattach u under v
           old := parent[u]
           // update old parent capacity
           oldCap := 2 - childCnt[old]
           childCnt[old]--
           newCapOld := 2 - childCnt[old]
           if old != 0 {
               if oldCap == 0 && newCapOld > 0 {
                   avail[depthArr[old]] = append(avail[depthArr[old]], old)
               } else if oldCap > 0 && newCapOld == 0 {
                   // remove old from avail
                   d0 := depthArr[old]
                   for i, x := range avail[d0] {
                       if x == old {
                           avail[d0] = append(avail[d0][:i], avail[d0][i+1:]...)
                           break
                       }
                   }
               }
           }
           // update new parent capacity
           pCapOld := 2 - childCnt[v]
           childCnt[v]++
           pCapNew := 2 - childCnt[v]
           if pCapOld > 0 && pCapNew == 0 {
               // remove v
               for i, x := range avail[dv] {
                   if x == v {
                       avail[dv] = append(avail[dv][:i], avail[dv][i+1:]...)
                       break
                   }
               }
           }
           // assign
           parent[u] = v
           depthArr[u] = nd
           // u remains leaf
           leaves[nd] = append(leaves[nd], u)
           if nd > leafMax {
               leafMax = nd
           }
       }
       // output
       fmt.Fprintln(out, "YES")
       for i := 2; i <= n; i++ {
           fmt.Fprint(out, parent[i])
           if i < n {
               fmt.Fprint(out, " ")
           }
       }
       fmt.Fprintln(out)
   }
}
