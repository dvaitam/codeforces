package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n uint64
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // l0 = L(0), l1 = L(1) in leaf-minimum sequence: L(h) = L(h-1) + L(h-2)
   var l0, l1 uint64 = 1, 2
   // h corresponds to current index for l1 (initial h=1 for l1)
   h := uint64(1)
   // increase h while l1 <= n
   for l1 <= n {
       h++
       l0, l1 = l1, l0+l1
   }
   // number of games the winner can play is h-1
   // (since when l1 > n, the last valid height is h-1)
   fmt.Println(h - 1)
}
