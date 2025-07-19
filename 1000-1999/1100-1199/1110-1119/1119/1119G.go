package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var N, M int
   fmt.Fscan(reader, &N, &M)
   H := make([]int64, M)
   for i := 0; i < M; i++ {
       fmt.Fscan(reader, &H[i])
   }
   // compute cutoffs
   cutoffs := make([]int64, 0, M+1)
   cutoffs = append(cutoffs, 0)
   cutoffs = append(cutoffs, int64(N))
   var cur int64
   for i := 0; i+1 < M; i++ {
       cur += H[i]
       cutoffs = append(cutoffs, cur%int64(N))
   }
   sort.Slice(cutoffs, func(i, j int) bool { return cutoffs[i] < cutoffs[j] })
   sizes := make([]int64, M)
   for i := 0; i < M; i++ {
       sizes[i] = cutoffs[i+1] - cutoffs[i]
   }
   // build answer sequence
   ans := make([]int, 0)
   ind := 0
   for i := 0; i < M; i++ {
       v := H[i]
       for v > 0 {
           ans = append(ans, i)
           v -= sizes[ind%M]
           ind++
       }
   }
   // pad to multiple of M
   for len(ans)%M != 0 {
       ans = append(ans, 0)
   }
   // output
   fmt.Fprintln(writer, len(ans)/M)
   // sizes line
   for i := 0; i < M; i++ {
       if i+1 == M {
           fmt.Fprintln(writer, sizes[i])
       } else {
           fmt.Fprint(writer, sizes[i], " ")
       }
   }
   // ans lines
   for i, x := range ans {
       if (i+1)%M == 0 {
           fmt.Fprintln(writer, x+1)
       } else {
           fmt.Fprint(writer, x+1, " ")
       }
   }
}
