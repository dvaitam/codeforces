package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   names := make([]string, n)
   idx := make(map[string]int, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       names[i] = s
       idx[s] = i
   }
   G := make([]int, n)
   for i := 0; i < m; i++ {
       var a, b string
       fmt.Fscan(reader, &a, &b)
       ia, ib := idx[a], idx[b]
       G[ia] |= 1 << ib
       G[ib] |= 1 << ia
   }
   var ans, bestMask int
   total := 1 << n
   for mask := 0; mask < total; mask++ {
       cur := mask
       for j := 0; j < n; j++ {
           if mask&(1<<j) != 0 {
               cur &^= G[j]
           }
       }
       cnt := bits.OnesCount(uint(cur))
       if cnt > ans {
           ans = cnt
           bestMask = cur
       }
   }
   var result []string
   for j := 0; j < n; j++ {
       if bestMask&(1<<j) != 0 {
           result = append(result, names[j])
       }
   }
   if len(result) == 0 {
       // fallback to first
       result = []string{names[0]}
   } else {
       sort.Strings(result)
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, len(result))
   for _, s := range result {
       fmt.Fprintln(writer, s)
   }
}
