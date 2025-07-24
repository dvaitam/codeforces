package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // If length is odd, cannot return to origin
   if n%2 == 1 {
       fmt.Println(-1)
       return
   }
   // Count directions
   l0, r0, u0, d0 := 0, 0, 0, 0
   for i := 0; i < n; i++ {
       switch s[i] {
       case 'L':
           l0++
       case 'R':
           r0++
       case 'U':
           u0++
       case 'D':
           d0++
       }
   }
   half := n / 2
   // best total deviation (sum of abs diffs)
   bestDev := 2 * n
   // Try k = target count for L and R
   for k := 0; k <= half; k++ {
       m := half - k
       dev := abs(l0-k) + abs(r0-k) + abs(u0-m) + abs(d0-m)
       if dev < bestDev {
           bestDev = dev
       }
   }
   // Each replacement fixes two deviations
   fmt.Println(bestDev / 2)
}
