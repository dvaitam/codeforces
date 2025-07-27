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
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   sum := int64(0)
   ans := 0
   seen := make(map[int64]struct{}, n*2)
   seen[0] = struct{}{}

   for i := 0; i < n; i++ {
       sum += a[i]
       if _, exists := seen[sum]; exists {
           // need to insert break before a[i]
           ans++
           // reset for new segment
           seen = make(map[int64]struct{}, n*2)
           // initial prefix sum zero for new segment
           seen[0] = struct{}{}
           // start sum from current element
           sum = a[i]
           seen[sum] = struct{}{}
       } else {
           seen[sum] = struct{}{}
       }
   }
   fmt.Fprintln(writer, ans)
}
