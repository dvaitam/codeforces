package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   arr := make([]int64, n+1)
   c := make([]int64, n)
   for i := 0; i < n; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       arr[i] = x
   }
   for i := 0; i < n; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       c[i] = x
   }
   for i := n - 1; i >= 0; i-- {
       mul := arr[i] * 100
       if mul > c[i] {
           arr[i] = c[i]
       } else {
           arr[i] = mul
           if arr[i+1] > arr[i] {
               if c[i] < arr[i+1] {
                   arr[i] = c[i]
               } else {
                   arr[i] = arr[i+1]
               }
           }
       }
   }
   slice := arr[:n]
   sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
   ans := float64(slice[0])
   tot := 1.0
   limit := n - k
   for i := 0; i < limit; i++ {
       tot *= float64(limit-i) / float64(n-i)
       if tot < 1e-18 {
           break
       }
       ans += tot * float64(slice[i+1]-slice[i])
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintf(writer, "%f\n", ans)
}
