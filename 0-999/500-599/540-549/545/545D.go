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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   times := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &times[i])
   }
   sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
   var sum int64
   var count int
   for _, t := range times {
       if sum <= t {
           count++
           sum += t
       }
   }
   fmt.Fprintln(writer, count)
}
