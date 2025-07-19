package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var tt int
   fmt.Fscan(in, &tt)
   for ; tt > 0; tt-- {
       var H, W int
       fmt.Fscan(in, &H, &W)
       for i := 0; i < H; i++ {
           for j := 0; j < W; j++ {
               var val byte
               // corners
               if (i == 0 || i == H-1) && (j == 0 || j == W-1) {
                   val = '1'
               } else if i == 0 || i == H-1 {
                   // top or bottom edges
                   if j == 1 || j == W-2 {
                       val = '0'
                   } else if j%2 == 0 {
                       val = '1'
                   } else {
                       val = '0'
                   }
               } else if j == 0 || j == W-1 {
                   // left or right edges
                   if i == 1 || i == H-2 {
                       val = '0'
                   } else if i%2 == 0 {
                       val = '1'
                   } else {
                       val = '0'
                   }
               } else {
                   val = '0'
               }
               out.WriteByte(val)
           }
           out.WriteByte('\n')
       }
       out.WriteByte('\n')
   }
}
