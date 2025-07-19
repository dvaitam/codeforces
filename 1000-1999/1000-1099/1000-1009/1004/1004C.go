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
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   // num[i] = number of distinct elements in a[i:]
   num := make([]int, n+1)
   seen := make(map[int]bool)
   for i := n - 1; i >= 0; i-- {
       if !seen[a[i]] {
           num[i] = num[i+1] + 1
           seen[a[i]] = true
       } else {
           num[i] = num[i+1]
       }
   }

   // count pairs (i<j) with distinct values
   seen = make(map[int]bool)
   var ans int64
   for i := 0; i < n-1; i++ {
       if !seen[a[i]] {
           ans += int64(num[i+1])
           seen[a[i]] = true
       }
   }

   fmt.Fprint(writer, ans)
}
