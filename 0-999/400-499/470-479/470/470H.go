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
   arr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }
   sort.Ints(arr)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i, v := range arr {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
}
