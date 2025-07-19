package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   var k int
   fmt.Fscan(reader, &s)
   fmt.Fscan(reader, &k)
   n := len(s)
   f := make([]int, 256)
   for i := 0; i < n; i++ {
       f[s[i]]++
   }
   p := make([]int, 256)
   m := n - k
   ans := 0
   for m > 0 {
       maxFreq := 0
       x := 0
       for i := 0; i < 256; i++ {
           if f[i] > maxFreq {
               maxFreq = f[i]
               x = i
           }
       }
       if maxFreq == 0 {
           break
       }
       p[x] = f[x]
       m -= p[x]
       ans++
       f[x] = 0
   }
   fmt.Println(ans)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < n; i++ {
       c := s[i]
       if p[c] > 0 {
           writer.WriteByte(c)
       }
   }
   writer.WriteByte('\n')
}
