package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   a := make([][]bool, n)
   for i := 0; i < n; i++ {
       a[i] = make([]bool, n)
       for j := 0; j < n; j++ {
           var x int
           fmt.Fscan(reader, &x)
           a[i][j] = x != 0
       }
   }
   // initialize random seed
   rand.Seed(time.Now().UnixNano())
   // prepare nodes
   tmp := make([]int, n)
   for i := 0; i < n; i++ {
       tmp[i] = i
   }
   bestTot := 0
   var bestAns []int
   const T = 7000
   for tIter := 0; tIter < T; tIter++ {
       rand.Shuffle(n, func(i, j int) { tmp[i], tmp[j] = tmp[j], tmp[i] })
       rec := make([]int, 0, n)
       for _, u := range tmp {
           ok := true
           for _, v := range rec {
               if !a[v][u] {
                   ok = false
                   break
               }
           }
           if ok {
               rec = append(rec, u)
           }
       }
       if len(rec) > bestTot {
           bestTot = len(rec)
           bestAns = append(([]int)(nil), rec...)
       }
   }
   // compute probabilities l and result
   l := make([]float64, n)
   if bestTot > 0 {
       p := float64(k) / float64(bestTot)
       for _, u := range bestAns {
           l[u] = p
       }
   }
   var ret float64
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           ret += l[i] * l[j]
       }
   }
   fmt.Printf("%.10f\n", ret)
}
