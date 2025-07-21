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
   fmt.Fscan(reader, &n)
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // prefix sum of positive values
   posPref := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       posPref[i] = posPref[i-1]
       if a[i] > 0 {
           posPref[i] += a[i]
       }
   }
   // map value to positions
   posMap := make(map[int64][]int)
   for i := 1; i <= n; i++ {
       v := a[i]
       posMap[v] = append(posMap[v], i)
   }
   // find best pair
   var bestSum int64 = -1<<62
   var bestL, bestR int
   for v, vec := range posMap {
       if len(vec) < 2 {
           continue
       }
       l := vec[0]
       r := vec[len(vec)-1]
       // sum = a[l] + a[r] + sum of positives between
       midSum := int64(0)
       if r > l+1 {
           midSum = posPref[r-1] - posPref[l]
       }
       total := v + v + midSum
       if total > bestSum {
           bestSum = total
           bestL = l
           bestR = r
       }
   }
   // mark kept
   keep := make([]bool, n+1)
   keep[bestL] = true
   keep[bestR] = true
   for i := bestL + 1; i < bestR; i++ {
       if a[i] > 0 {
           keep[i] = true
       }
   }
   // collect removed
   removed := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if !keep[i] {
           removed = append(removed, i)
       }
   }
   // output
   fmt.Fprintln(writer, bestSum, len(removed))
   for i, idx := range removed {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, idx)
   }
   if len(removed) > 0 {
       fmt.Fprintln(writer)
   }
}
