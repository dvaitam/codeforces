package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   var S int
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       S += a[i]
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &b[i])
   }
   // sort capacities descending
   caps := make([]int, n)
   copy(caps, b)
   sort.Sort(sort.Reverse(sort.IntSlice(caps)))
   // find minimal k: sum of top k caps >= S
   sum := 0
   k := 0
   for i, v := range caps {
       sum += v
       if sum >= S {
           k = i + 1
           break
       }
   }
   // capacity threshold
   capK := caps[k-1]
   // sum ai for bottles with b > capK
   var sumAiGreater int
   var cntGreater int
   var eqA []int
   for i := 0; i < n; i++ {
       if b[i] > capK {
           cntGreater++
           sumAiGreater += a[i]
       } else if b[i] == capK {
           eqA = append(eqA, a[i])
       }
   }
   // select remaining eq bottles with largest a
   need := k - cntGreater
   sort.Sort(sort.Reverse(sort.IntSlice(eqA)))
   sumEq := 0
   if need > len(eqA) {
       need = len(eqA)
   }
   for i := 0; i < need; i++ {
       sumEq += eqA[i]
   }
   // total kept ai
   kept := sumAiGreater + sumEq
   // time is amount moved from others
   t := S - kept
   // output k and t
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintf(w, "%d %d", k, t)
}
