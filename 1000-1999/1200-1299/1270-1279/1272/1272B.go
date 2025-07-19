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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var s string
       fmt.Fscan(reader, &s)
       var cnt [256]int
       for _, c := range s {
           cnt[c]++
       }
       horiz := min(cnt['L'], cnt['R'])
       vert := min(cnt['U'], cnt['D'])

       // No moves possible
       if horiz == 0 && vert == 0 {
           fmt.Fprintln(writer, 0)
           continue
       }
       // Only vertical moves
       if horiz == 0 {
           fmt.Fprintln(writer, 2)
           fmt.Fprintln(writer, "UD")
           continue
       }
       // Only horizontal moves
       if vert == 0 {
           fmt.Fprintln(writer, 2)
           fmt.Fprintln(writer, "LR")
           continue
       }
       // General case: rectangle path
       total := (horiz + vert) * 2
       fmt.Fprintln(writer, total)
       // Build path: U^vert, R^horiz, D^vert, L^horiz
       path := make([]byte, 0, total)
       for j := 0; j < vert; j++ {
           path = append(path, 'U')
       }
       for j := 0; j < horiz; j++ {
           path = append(path, 'R')
       }
       for j := 0; j < vert; j++ {
           path = append(path, 'D')
       }
       for j := 0; j < horiz; j++ {
           path = append(path, 'L')
       }
       fmt.Fprintln(writer, string(path))
   }
}
