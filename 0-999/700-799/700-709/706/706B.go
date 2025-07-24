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

   var n int
   fmt.Fscan(reader, &n)
   prices := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &prices[i])
   }
   sort.Ints(prices)

   var q int
   fmt.Fscan(reader, &q)
   for i := 0; i < q; i++ {
       var m int
       fmt.Fscan(reader, &m)
       // count shops with price <= m
       cnt := sort.Search(len(prices), func(j int) bool { return prices[j] > m })
       fmt.Fprintln(writer, cnt)
   }
}
