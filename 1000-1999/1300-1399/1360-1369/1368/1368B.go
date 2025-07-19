package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var k int64
   fmt.Fscan(in, &k)
   s := "codeforces"
   for l := 10; ; l++ {
       // Distribute l occurrences among 10 characters
       c := make([]int, 10)
       j := l
       for i := 0; i < 10; i++ {
           c[i] = j / (10 - i)
           j -= c[i]
       }
       // Compute product of counts
       var res int64 = 1
       for i := 0; i < 10; i++ {
           res *= int64(c[i])
       }
       if res < k {
           continue
       }
       // Output result
       for i := 0; i < 10; i++ {
           for rep := 0; rep < c[i]; rep++ {
               out.WriteByte(s[i])
           }
       }
       out.WriteByte('\n')
       break
   }
}
