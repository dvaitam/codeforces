package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func sign(x int) int {
   if x < 0 {
       return -1
   } else if x > 0 {
       return 1
   }
   return 0
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var kx, ky int
   // read king position
   if _, err := fmt.Fscan(reader, &kx, &ky); err != nil {
       return
   }
   // read rooks
   const R = 666
   rx := make([]int, R)
   ry := make([]int, R)
   for i := 0; i < R; i++ {
       fmt.Fscan(reader, &rx[i], &ry[i])
   }

   // target center
   const CX = 500
   const CY = 500

   for moves := 0; moves < 2000; moves++ {
       // compute next king move towards center
       dx := sign(CX - kx)
       dy := sign(CY - ky)
       nx := kx + dx
       ny := ky + dy
       // make move
       fmt.Fprintf(writer, "%d %d\n", nx, ny)
       writer.Flush()
       kx, ky = nx, ny

       // check if any rook aligns
       for i := 0; i < R; i++ {
           if rx[i] == kx || ry[i] == ky {
               // king is in check -> win, exit
               return
           }
       }

       // read rook move
       var idx, x, y int
       if _, err := fmt.Fscan(reader, &idx, &x, &y); err != nil {
           return
       }
       idx-- // zero-based
       if idx >= 0 && idx < R {
           rx[idx] = x
           ry[idx] = y
       }
   }
}
