package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var h int
   var n int64
   if _, err := fmt.Fscan(reader, &h, &n); err != nil {
       return
   }
   var ans int64 = 0
   // dir: 0 means left is prioritized next, 1 means right
   var dir int = 0
   for i := h; i > 0; i-- {
       // size of half leaves at this level
       half := int64(1) << (i - 1)
       if (n <= half && dir == 0) || (n > half && dir == 1) {
           // go directly
           ans += 1
           // flip priority
           dir ^= 1
       } else {
           // traverse the entire other subtree
           ans += int64(1) << i
           // dir remains
       }
       // move to next subtree
       if n > half {
           n -= half
       }
   }
   fmt.Println(ans)
}
