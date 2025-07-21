package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // next occurrence of "bear" starting at or after index i
   nextOcc := make([]int, n)
   const inf = int(1e9)
   last := inf
   for i := n - 1; i >= 0; i-- {
       if i <= n-4 && s[i] == 'b' && s[i+1] == 'e' && s[i+2] == 'a' && s[i+3] == 'r' {
           last = i
       }
       nextOcc[i] = last
   }
   var ans int64
   for i := 0; i < n; i++ {
       p := nextOcc[i]
       // valid start position must allow full "bear" (4 chars)
       if p <= n-4 {
           // substrings from i..j with j >= p+3
           ans += int64(n - (p + 3))
       }
   }
   fmt.Println(ans)
}
