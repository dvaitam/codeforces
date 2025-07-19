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

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n, k int
       fmt.Fscan(reader, &n, &k)
       // n: number of elements, k: max value
       s := make([]int, 0, n)
       present := make(map[int]bool, n)
       // start with [1, 2]
       if n >= 1 {
           s = append(s, 1)
           present[1] = true
       }
       if n >= 2 {
           s = append(s, 2)
           present[2] = true
       }
       cnt := 2
       // build with increasing gaps
       for i := 4; i <= k && len(s) < n; {
           s = append(s, i)
           present[i] = true
           cnt++
           i += cnt
       }
       // fill remaining from k downwards
       if len(s) < n {
           for x := k; x >= 1 && len(s) < n; x-- {
               if !present[x] {
                   s = append(s, x)
               }
           }
       }
       sort.Ints(s)
       // output
       for _, v := range s {
           fmt.Fprintf(writer, "%d ", v)
       }
       fmt.Fprintln(writer)
   }
}
