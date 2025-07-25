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

   var n, m, ta, tb, k int
   fmt.Fscan(reader, &n, &m, &ta, &tb, &k)
   a := make([]int, n)
   b := make([]int, m)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &b[i])
   }
   // If we can cancel all flights from A or B, connection impossible
   if k >= n || k >= m {
       fmt.Fprintln(writer, -1)
       return
   }
   ans := 0
   // Try removing x flights from A and k-x flights from B
   for x := 0; x <= k; x++ {
       y := k - x
       // Choose the (x)-th flight from A (0-based) after cancellations
       departTime := a[x] + ta
       // Find earliest B->C flight we can catch
       j := sort.Search(m, func(i int) bool { return b[i] >= departTime })
       // After removing y flights among those we could catch, index shifts by y
       if j+y >= m {
           fmt.Fprintln(writer, -1)
           return
       }
       arrival := b[j+y] + tb
       if arrival > ans {
           ans = arrival
       }
   }
   fmt.Fprintln(writer, ans)
}
