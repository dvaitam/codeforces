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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   const mod = 1000000007
   // primes up to 2e6, mark seen
   maxP := 2000000
   seen := make([]bool, maxP+1)
   ans := 1
   for i := 0; i < n; i++ {
       var p int
       fmt.Fscan(reader, &p)
       if !seen[p] {
           seen[p] = true
           ans = int((int64(ans) * int64(p)) % mod)
       }
   }
   fmt.Fprint(writer, ans)
}
