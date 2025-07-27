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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   // positions maps value to its list of positions
   positions := make(map[int][]int)
   uniq := make([]int, 0)
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       uniq = uniq[:0]
       for i := 1; i <= n; i++ {
           var x int
           fmt.Fscan(reader, &x)
           if _, ok := positions[x]; !ok {
               uniq = append(uniq, x)
           }
           positions[x] = append(positions[x], i)
       }
       ans := n
       for _, x := range uniq {
           pos := positions[x]
           cnt := 0
           // segment before first occurrence
           if pos[0] != 1 {
               cnt++
           }
           // gaps between occurrences
           for i := 1; i < len(pos); i++ {
               if pos[i] != pos[i-1]+1 {
                   cnt++
               }
           }
           // segment after last occurrence
           if pos[len(pos)-1] != n {
               cnt++
           }
           if cnt < ans {
               ans = cnt
           }
           delete(positions, x)
       }
       fmt.Fprintln(writer, ans)
   }
}
