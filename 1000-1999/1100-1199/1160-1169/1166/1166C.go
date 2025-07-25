package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       if x < 0 {
           x = -x
       }
       b[i] = x
   }
   sort.Slice(b, func(i, j int) bool {
       return b[i] < b[j]
   })
   var ans int64
   j := 0
   for i := 0; i < n; i++ {
       if j < i {
           j = i
       }
       for j+1 < n && b[j+1] <= 2*b[i] {
           j++
       }
       ans += int64(j - i)
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
