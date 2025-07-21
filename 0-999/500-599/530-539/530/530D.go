package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   removed := make([]bool, 1001)
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       if a < 1 {
           a = 1
       }
       if b > 1000 {
           b = 1000
       }
       for j := a; j <= b; j++ {
           removed[j] = true
       }
   }
   var result []int
   for i := 1; i <= 1000; i++ {
       if !removed[i] {
           result = append(result, i)
       }
   }
   // Output count and elements
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, len(result))
   for _, v := range result {
       fmt.Fprint(writer, " ", v)
   }
   fmt.Fprintln(writer)
}
