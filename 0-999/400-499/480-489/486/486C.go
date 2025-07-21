package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, p int
   if _, err := fmt.Fscan(reader, &n, &p); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)

   // positions are 1-based
   half := n / 2
   // reflect cursor to left half for minimal movement
   if p > half {
       p = n - p + 1
   }

   totalChange := 0
   left, right := n, 0
   // compute character change costs and record bounds
   for i := 1; i <= half; i++ {
       a := s[i-1]
       b := s[n-i]
       // compute absolute difference in alphabet positions
       diff := int(a) - int(b)
       if diff < 0 {
           diff = -diff
       }
       // minimal cyclic distance in alphabet
       cost := min(diff, 26-diff)
       if cost > 0 {
           totalChange += cost
           if i < left {
               left = i
           }
           if i > right {
               right = i
           }
       }
   }
   if totalChange == 0 {
       fmt.Fprint(writer, 0)
       return
   }
   // compute movement cost
   mv := 0
   // if only one position, left==right
   // movement cost is distance from p to that pos
   // general: cover from left to right
   span := right - left
   // movement: go to closer end, then traverse span
   mv = span + min(abs(p-left), abs(p-right))

   // output total operations
   fmt.Fprint(writer, totalChange+mv)
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}
