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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   arr := make([]int64, n)
   var sum int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
       sum += arr[i]
   }
   // find largest and second largest
   var max1, max2 int64
   var idx1, idx2 int
   if n > 0 {
       max1 = arr[0]
       idx1 = 0
   }
   for i := 1; i < n; i++ {
       if arr[i] > max1 {
           max1 = arr[i]
           idx1 = i
       }
   }
   // compute sum of others after removing max1
   l := sum - max1
   flag := make([]bool, n)
   res := make([]int, 0)
   // check for all except idx1
   for i := 0; i < n; i++ {
       if i != idx1 && l-arr[i] == max1 {
           flag[i] = true
           res = append(res, i+1)
       }
   }
   // find second max among remaining
   for i := 0; i < n; i++ {
       if i == idx1 {
           continue
       }
       if arr[i] > max2 {
           max2 = arr[i]
           idx2 = i
       }
   }
   // check idx1 if qualifies
   if sum-max1-max2 == max2 && !flag[idx1] {
       res = append(res, idx1+1)
   }
   // output
   fmt.Fprintln(writer, len(res))
   for i, v := range res {
       if i > 0 {
           writer.WriteString(" ")
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
