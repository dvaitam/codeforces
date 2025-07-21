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
   // read the base string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   var k int64
   fmt.Fscan(reader, &k)
   // read weights for 'a' to 'z'
   weights := make([]int64, 26)
   var maxW int64
   for i := 0; i < 26; i++ {
       fmt.Fscan(reader, &weights[i])
       if weights[i] > maxW {
           maxW = weights[i]
       }
   }

   var total int64
   // sum for original string
   for i, ch := range s {
       w := weights[ch-'a']
       total += w * int64(i+1)
   }
   n := int64(len(s))
   // sum for k inserted letters of max weight at the end
   // positions: n+1 .. n+k
   // sum of positions = k*n + k*(k+1)/2
   total += maxW * (k*n + k*(k+1)/2)

   fmt.Fprint(writer, total)
}
