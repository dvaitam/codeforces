package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   lengths := make([]int, n)
   var s string
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &s)
       lengths[i] = len(s)
   }
   var target string
   fmt.Fscan(reader, &target)
   L := len(target)
   less, equal := 0, 0
   for _, l := range lengths {
       if l < L {
           less++
       } else if l == L {
           equal++
       }
   }
   bestPos := less + 1
   worstPos := less + equal
   bestTime := bestPos + ((bestPos-1)/k)*5
   worstTime := worstPos + ((worstPos-1)/k)*5
   fmt.Println(bestTime, worstTime)
}
