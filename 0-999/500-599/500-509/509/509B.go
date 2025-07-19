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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   arr := make([]int, n)
   mini := int(1e9)
   maxi := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
       if arr[i] < mini {
           mini = arr[i]
       }
       if arr[i] > maxi {
           maxi = arr[i]
       }
   }
   if maxi-mini > k {
       fmt.Fprint(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   for _, a := range arr {
       pop := 1
       flag := false
       for j := 1; j <= a; j++ {
           // print current color
           fmt.Fprint(writer, pop, " ")
           if flag {
               pop++
           }
           if j == mini && !flag {
               flag = true
           }
       }
       fmt.Fprintln(writer)
   }
}
