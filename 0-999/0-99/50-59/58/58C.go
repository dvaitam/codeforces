package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   freq := make(map[int]int)
   maxf := 0
   for i, v := range a {
       // distance from nearest edge
       d := i
       if n-1-i < d {
           d = n - 1 - i
       }
       x := v - d
       if x > 0 {
           freq[x]++
           if freq[x] > maxf {
               maxf = freq[x]
           }
       }
   }
   // minimal changes = total trees - maximal matches
   fmt.Println(n - maxf)
}
