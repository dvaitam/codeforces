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
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       if n < 3 {
           fmt.Fprintln(writer, 0)
           continue
       }
       sort.Ints(a)
       var res int64
       l := 0
       for r := 0; r < n; r++ {
           for l < r && a[r]-a[l] > 2 {
               l++
           }
           cnt := r - l
           if cnt >= 2 {
               // choose any two from cnt to pair with r
               res += int64(cnt) * int64(cnt-1) / 2
           }
       }
       fmt.Fprintln(writer, res)
   }
}
