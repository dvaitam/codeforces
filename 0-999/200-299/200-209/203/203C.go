package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type customer struct {
   idx   int
   space int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var d, a, b int64
   if _, err := fmt.Fscan(reader, &n, &d); err != nil {
       return
   }
   fmt.Fscan(reader, &a, &b)

   arr := make([]customer, n)
   for i := 0; i < n; i++ {
       var x, y int64
       fmt.Fscan(reader, &x, &y)
       arr[i] = customer{idx: i + 1, space: a*x + b*y}
   }
   sort.Slice(arr, func(i, j int) bool {
       return arr[i].space < arr[j].space
   })

   var sum int64
   var count int
   for count < n && sum+arr[count].space <= d {
       sum += arr[count].space
       count++
   }

   fmt.Fprintln(writer, count)
   if count > 0 {
       for i := 0; i < count; i++ {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, arr[i].idx)
       }
       fmt.Fprintln(writer)
   }
}
