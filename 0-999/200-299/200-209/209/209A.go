package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // f0: count of zebroid subsequences ending at even position
   // f1: count of zebroid subsequences ending at odd position
   var f0, f1 int64
   for i := 1; i <= n; i++ {
       var add int64 = 1
       if i%2 == 0 {
           // even position extends subsequences ending at odd
           add = (add + f1) % mod
           f0 = (f0 + add) % mod
       } else {
           // odd position extends subsequences ending at even
           add = (add + f0) % mod
           f1 = (f1 + add) % mod
       }
   }
   ans := (f0 + f1) % mod
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, ans)
}
