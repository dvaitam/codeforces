package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   var s string
   fmt.Fscan(reader, &s)
   b := []byte(s)
   V := make([][]int, 0)
   tot := 0

   // Build sequence of moves
   for {
       cur := 0
       mn := n + 1
       v := make([]int, 0)
       for i := 1; i < n; i++ {
           if b[i] == 'L' {
               cur++
           } else {
               cur--
           }
           if b[i] == 'L' && b[i-1] == 'R' {
               if cur < mn {
                   mn = cur
                   v = v[:0]
                   v = append(v, i)
               } else if cur == mn {
                   v = append(v, i)
               }
           }
       }
       if len(v) > 0 {
           // apply swaps
           for _, x := range v {
               b[x-1], b[x] = b[x], b[x-1]
           }
           V = append(V, v)
           tot += len(v)
       } else {
           break
       }
   }

   // check feasibility
   if len(V) > k || tot < k {
       fmt.Fprintln(writer, -1)
       return
   }

   // output exactly k operations
   y := len(V)
   for _, v := range V {
       y--
       if len(v)+y <= k {
           for _, x := range v {
               fmt.Fprintln(writer, 1, x)
           }
           y += len(v)
       } else {
           cnt := k - y - 1
           for i := 0; i < cnt; i++ {
               fmt.Fprintln(writer, 1, v[i])
               y++
           }
           // remaining moves in one line
           fmt.Fprint(writer, len(v)-cnt)
           for i := cnt; i < len(v); i++ {
               fmt.Fprint(writer, " ", v[i])
           }
           fmt.Fprint(writer, "\n")
           y++
       }
   }
}
