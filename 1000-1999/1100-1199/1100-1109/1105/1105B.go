package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   var s string
   fmt.Fscan(reader, &n, &k, &s)
   res := make([]int, 26)
   i := 0
   for i < n {
       x := s[i]
       c := 0
       for i < n && s[i] == x {
           c++
           i++
       }
       res[x-'a'] += c / k
   }
   mx := 0
   for _, v := range res {
       if v > mx {
           mx = v
       }
   }
   fmt.Println(mx)
}
