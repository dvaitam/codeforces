package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct {
   val int64
   idx int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   var k int64
   fmt.Fscan(in, &n, &m, &k)
   // deques for maxima
   deques := make([][]pair, m)
   currentMax := make([]int64, m)
   bestMax := make([]int64, m)
   l := 0
   bestLen := 0

   for r := 0; r < n; r++ {
       // read row
       a := make([]int64, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(in, &a[j])
       }
       // add to deques
       for j := 0; j < m; j++ {
           dq := &deques[j]
           // pop smaller or equal
           for len(*dq) > 0 && (*dq)[len(*dq)-1].val <= a[j] {
               *dq = (*dq)[:len(*dq)-1]
           }
           *dq = append(*dq, pair{val: a[j], idx: r})
       }
       // compute current sum of maxima
       var sum int64
       for j := 0; j < m; j++ {
           if len(deques[j]) > 0 {
               currentMax[j] = deques[j][0].val
           } else {
               currentMax[j] = 0
           }
           sum += currentMax[j]
       }
       // shrink window from left if over budget
       for sum > k && l <= r {
           for j := 0; j < m; j++ {
               dq := &deques[j]
               if len(*dq) > 0 && (*dq)[0].idx == l {
                   *dq = (*dq)[1:]
               }
           }
           l++
           // recompute sum after pop
           sum = 0
           for j := 0; j < m; j++ {
               if len(deques[j]) > 0 {
                   currentMax[j] = deques[j][0].val
               } else {
                   currentMax[j] = 0
               }
               sum += currentMax[j]
           }
       }
       // update best
       curLen := r - l + 1
       if curLen > bestLen {
           bestLen = curLen
           copy(bestMax, currentMax)
       }
   }
   // output best shots per type
   for j := 0; j < m; j++ {
       if j > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, bestMax[j])
   }
   out.WriteByte('\n')
}
