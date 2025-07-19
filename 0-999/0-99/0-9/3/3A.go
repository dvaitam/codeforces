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

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func sign(x int) int {
   if x < 0 {
       return -1
   }
   if x > 0 {
       return 1
   }
   return 0
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var s, t string
   if _, err := fmt.Fscan(reader, &s, &t); err != nil {
       return
   }

   x1 := int(s[0] - 'a')
   y1 := int(s[1] - '1')
   x2 := int(t[0] - 'a')
   y2 := int(t[1] - '1')

   steps := max(abs(x1-x2), abs(y1-y2))
   fmt.Fprintln(writer, steps)

   for x1 != x2 || y1 != y2 {
       dx := sign(x2 - x1)
       dy := sign(y2 - y1)
       dir := ""
       if dx < 0 {
           dir += "L"
       } else if dx > 0 {
           dir += "R"
       }
       if dy > 0 {
           dir += "U"
       } else if dy < 0 {
           dir += "D"
       }
       fmt.Fprintln(writer, dir)
       x1 += dx
       y1 += dy
   }
}
