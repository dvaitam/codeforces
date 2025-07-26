package main

import (
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   var n int
   if _, err := fmt.Fscan(os.Stdin, &n); err != nil {
       return
   }
   res := 1
   for i := 2; i <= n; i++ {
       res = res * i % mod
   }
   fmt.Println(res)
}
