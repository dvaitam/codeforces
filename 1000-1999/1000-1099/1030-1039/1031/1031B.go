package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n-1)
   b := make([]int, n-1)
   for i := 0; i < n-1; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 0; i < n-1; i++ {
       fmt.Fscan(reader, &b[i])
   }

   // dpPrev[v] = reachable with t[i]=v
   dpPrev := [4]bool{}
   dpCurr := [4]bool{}
   // parent[i][w] = previous v for t[i]=w
   parent := make([][4]int, n)
   // initialize parents to -1
   for i := 0; i < n; i++ {
       for v := 0; v < 4; v++ {
           parent[i][v] = -1
       }
   }
   // initial: position 0 (t1) can be any
   for v := 0; v < 4; v++ {
       dpPrev[v] = true
       // parent[0][v] remains -1
   }

   // build dp
   for i := 0; i < n-1; i++ {
       for v := 0; v < 4; v++ {
           dpCurr[v] = false
       }
       for v := 0; v < 4; v++ {
           if !dpPrev[v] {
               continue
           }
           for w := 0; w < 4; w++ {
               if (v|w) == a[i] && (v&w) == b[i] {
                   if !dpCurr[w] {
                       dpCurr[w] = true
                       parent[i+1][w] = v
                   }
               }
           }
       }
       dpPrev = dpCurr
   }

   // find any reachable at end
   endVal := -1
   for v := 0; v < 4; v++ {
       if dpPrev[v] {
           endVal = v
           break
       }
   }
   if endVal < 0 {
       fmt.Fprintln(writer, "NO")
       return
   }

   // reconstruct t
   t := make([]int, n)
   t[n-1] = endVal
   for i := n - 1; i > 0; i-- {
       t[i-1] = parent[i][t[i]]
   }

   fmt.Fprintln(writer, "YES")
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, t[i])
   }
   fmt.Fprintln(writer)
}
