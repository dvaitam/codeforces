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

   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)
   ans := make([]int, m)
   for i := 1; i <= n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       for j := 1; j <= m; j++ {
           switch line[j-1] {
           case 'L':
               target := j - i + 1
               if target >= 1 && target <= m {
                   ans[target-1]++
               }
           case 'R':
               target := j + i - 1
               if target >= 1 && target <= m {
                   ans[target-1]++
               }
           case 'U':
               if i%2 == 1 {
                   ans[j-1]++
               }
           }
       }
   }
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(v))
   }
   writer.WriteByte('\n')
}
