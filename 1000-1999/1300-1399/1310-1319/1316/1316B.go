package main

import (
   "bufio"
   "fmt"
   "os"
)

func reverse(s string) string {
   b := []byte(s)
   for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
       b[i], b[j] = b[j], b[i]
   }
   return string(b)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       var s string
       fmt.Fscan(reader, &n, &s)
       // initialize with k=1
       bestK := 1
       var best string
       if n == 1 {
           best = s
       } else {
           // k=1: since k<n for n>1, t = s[0] + s[1:] + s[:0] = s
           best = s
       }
       // try k from 2 to n
       for k := 2; k <= n; k++ {
           var tstr string
           if k < n {
               // s[k-1] + s[k:] + s[:k-1]
               tstr = s[k-1:k] + s[k:n] + s[0:k-1]
           } else {
               // k == n: full reverse
               tstr = reverse(s)
           }
           if tstr < best {
               best = tstr
               bestK = k
           }
       }
       fmt.Fprintln(writer, best)
       fmt.Fprintln(writer, bestK)
   }
}
