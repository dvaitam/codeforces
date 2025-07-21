package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var recipe string
   fmt.Fscan(reader, &recipe)

   var nb, ns, nc int64
   fmt.Fscan(reader, &nb, &ns, &nc)
   var pb, ps, pc int64
   fmt.Fscan(reader, &pb, &ps, &pc)
   var r int64
   fmt.Fscan(reader, &r)

   // Count ingredients per hamburger
   var needB, needS, needC int64
   for _, ch := range recipe {
       switch ch {
       case 'B':
           needB++
       case 'S':
           needS++
       case 'C':
           needC++
       }
   }

   // Function to check if x hamburgers can be made
   can := func(x int64) bool {
       // Total cost to buy missing ingredients
       var cost int64
       // Bread
       req := needB*x - nb
       if req > 0 {
           cost += req * pb
       }
       // Sausage
       req = needS*x - ns
       if req > 0 {
           cost += req * ps
       }
       // Cheese
       req = needC*x - nc
       if req > 0 {
           cost += req * pc
       }
       return cost <= r
   }

   // Binary search maximum x
   var lo, hi int64 = 0, 20000000000000
   var ans int64
   for lo <= hi {
       mid := (lo + hi) / 2
       if can(mid) {
           ans = mid
           lo = mid + 1
       } else {
           hi = mid - 1
       }
   }
   fmt.Println(ans)
}
