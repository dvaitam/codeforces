package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   freq := make(map[int]int)
   for i := 0; i < m; i++ {
       var a int
       fmt.Fscan(in, &a)
       freq[a]++
   }
   if n > m {
       fmt.Println(0)
       return
   }
   // collect frequencies
   counts := make([]int, 0, len(freq))
   for _, v := range freq {
       if v > 0 {
           counts = append(counts, v)
       }
   }
   sort.Slice(counts, func(i, j int) bool { return counts[i] > counts[j] })

   // find max days x such that sum(counts[i]/x) >= n
   var ans int
   for days := 1; ; days++ {
       total := 0
       for _, v := range counts {
           total += v / days
       }
       if total < n {
           ans = days - 1
           break
       }
   }
   fmt.Println(ans)
}
