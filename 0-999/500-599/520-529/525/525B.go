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

   var s string
   // read initial string
   fmt.Fscan(reader, &s)
   n := len(s)
   // read number of operations
   var m int
   fmt.Fscan(reader, &m)
   // half length (floor)
   half := n / 2
   // frequency of operations starting at each ai
   freq := make([]int, half+2)
   for i := 0; i < m; i++ {
       var a int
       fmt.Fscan(reader, &a)
       if a <= half {
           freq[a]++
       }
   }
   // convert to mutable byte slice
   b := []byte(s)
   // apply swaps based on prefix sums
   sum := 0
   for i := 1; i <= half; i++ {
       sum += freq[i]
       if sum&1 == 1 {
           // swap positions i-1 and n-i
           b[i-1], b[n-i] = b[n-i], b[i-1]
       }
   }
   // output result
   writer.Write(b)
   writer.WriteByte('\n')
}
