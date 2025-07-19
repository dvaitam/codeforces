package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // Read n and k
   var n int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // Read array a[1..n]
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Sort ascending
   sort.Slice(a[1:], func(i, j int) bool { return a[i+1] < a[j+1] })

   // Generate and print sequences
   for i := 0; i <= n; i++ {
       // number of trailing fixed elements = n-i
       for j := 1; j <= n-i; j++ {
           if k <= 0 {
               return
           }
           k--
           // first element: i+1
           writer.WriteString(strconv.Itoa(i + 1))
           // print elements a[p] for p in n-i+1..n
           for p := n - i + 1; p <= n; p++ {
               writer.WriteByte(' ')
               writer.WriteString(strconv.Itoa(a[p]))
           }
           // then element a[j]
           writer.WriteByte(' ')
           writer.WriteString(strconv.Itoa(a[j]))
           writer.WriteByte('\n')
       }
   }
}
