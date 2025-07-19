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
       // Read values
       type pair struct{ val int64; idx int }
       arr := make([]pair, n)
       for i := 0; i < n; i++ {
           var v int64
           fmt.Fscan(reader, &v)
           arr[i] = pair{v, i + 1}
       }
       // Sort by value descending
       sort.Slice(arr, func(i, j int) bool {
           return arr[i].val > arr[j].val
       })
       // Prepare answers: positions from 0..n
       ans := make([]int64, n+1)
       var sum int64
       for i, p := range arr {
           k := int64(i/2 + 1)
           if i%2 == 0 {
               ans[p.idx] = k
           } else {
               ans[p.idx] = -k
           }
           sum += 2 * p.val * k
       }
       // Output total score
       fmt.Fprintln(writer, sum)
       // Output positions, including ans[0]=0
       // Print ans[0]
       writer.WriteByte('0')
       // Print ans[1..n]
       for i := 1; i <= n; i++ {
           writer.WriteByte(' ')
           fmt.Fprint(writer, ans[i])
       }
       writer.WriteByte('\n')
   }
}
