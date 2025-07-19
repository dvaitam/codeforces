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
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n+1)
       cnt := 0
       for i := 1; i <= n; i++ {
           fmt.Fscan(reader, &a[i])
           cnt += a[i]
       }
       b := make([]int, 0, cnt)
       for i := 1; i <= n; i++ {
           for j := 0; j < a[i]; j++ {
               b = append(b, i)
           }
       }
       // two pointers to pair smaller with larger
       start := 0
       x := cnt / 2
       // skip if equal to smallest
       for x < cnt && b[x] == b[0] {
           x++
       }
       var pairs [][2]int
       stx := -1
       // pair until pointers meet initial second start
       for x < cnt && start != stx {
           // skip equal elements
           for x < cnt-1 && b[start] == b[x] {
               x++
           }
           if x < cnt && b[start] != b[x] {
               pairs = append(pairs, [2]int{b[start], b[x]})
               if stx < 0 {
                   stx = x
               }
           }
           start++
           x++
       }
       // output
       fmt.Fprintln(writer, len(pairs))
       for _, p := range pairs {
           fmt.Fprintln(writer, p[0], p[1])
       }
   }
}
