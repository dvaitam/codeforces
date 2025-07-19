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

   var n int
   fmt.Fscan(reader, &n)
   s := make([]int, n)
   buck := make([]int, 510)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &s[i])
       buck[s[i]]++
   }
   calc := make([]int, 4)
   for i := range buck {
       if buck[i] <= 2 {
           calc[buck[i]]++
       } else {
           calc[3]++
       }
   }
   if calc[1]%2 == 1 && calc[3] == 0 {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   remains := calc[1] % 2
   half := calc[1] / 2
   t := 0
   res := make([]byte, n)
   for i := 0; i < n; i++ {
       cnt := buck[s[i]]
       switch {
       case cnt == 1:
           t++
           if t <= half {
               res[i] = 'A'
           } else {
               res[i] = 'B'
           }
       case cnt == 2:
           res[i] = 'A'
       default:
           if remains == 1 {
               remains = 0
               res[i] = 'A'
           } else {
               res[i] = 'B'
           }
       }
   }
   writer.Write(res)
   writer.WriteByte('\n')
}
