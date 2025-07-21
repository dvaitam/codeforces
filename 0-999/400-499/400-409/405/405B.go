package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   var s string
   fmt.Fscan(reader, &s)
   // Collect forces with sentinels
   type force struct { pos int; dir byte }
   forces := make([]force, 0, n+2)
   forces = append(forces, force{0, 'L'})
   for i := 1; i <= n; i++ {
       c := s[i-1]
       if c == 'L' || c == 'R' {
           forces = append(forces, force{i, c})
       }
   }
   forces = append(forces, force{n + 1, 'R'})

   // Count upright dominoes
   ans := 0
   for i := 0; i+1 < len(forces); i++ {
       left := forces[i]
       right := forces[i+1]
       d := right.pos - left.pos - 1
       if d <= 0 {
           continue
       }
       if left.dir == right.dir {
           // same direction, all fall
           continue
       }
       if left.dir == 'L' && right.dir == 'R' {
           // forces outward, all remain
           ans += d
       } else if left.dir == 'R' && right.dir == 'L' {
           // forces inward, middle remains if odd
           if d%2 == 1 {
               ans++
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
