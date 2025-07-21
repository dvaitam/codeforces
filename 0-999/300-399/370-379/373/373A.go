package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k int
   if _, err := fmt.Fscan(reader, &k); err != nil {
       return
   }
   // counts of panels per timing digit
   counts := make([]int, 10)
   for i := 0; i < 4; i++ {
       var line string
       fmt.Fscan(reader, &line)
       for _, c := range line {
           if c == '.' {
               continue
           }
           if c >= '1' && c <= '9' {
               counts[c-'0']++
           }
       }
   }
   // each time slot can be handled by two hands, each can press up to k panels
   limit := 2 * k
   for d := 1; d <= 9; d++ {
       if counts[d] > limit {
           fmt.Println("NO")
           return
       }
   }
   fmt.Println("YES")
}
