package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   const mod = 1000000009
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int64
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // For the given forest pattern, the number of suitable oriented routes
   // equals n^2 modulo mod
   ans := (n % mod) * (n % mod) % mod
   fmt.Fprintln(out, ans)
}
