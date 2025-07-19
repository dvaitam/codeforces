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
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       s := make([][]byte, n)
       cnt := make([]int, n)
       total := 0
       for i := 0; i < n; i++ {
           row := make([]byte, m)
           // read as string
           var str string
           fmt.Fscan(reader, &str)
           for j := 0; j < m; j++ {
               row[j] = str[j]
               if row[j] == '1' {
                   cnt[i]++
               }
           }
           total += cnt[i]
           s[i] = row
       }
       if total%n != 0 {
           fmt.Fprintln(writer, -1)
           continue
       }
       per := total / n
       // donors per column: rows with surplus and 1 in that column
       donors := make([][]int, m)
       deficit := make([]int, 0, n)
       for i := 0; i < n; i++ {
           if cnt[i] > per {
               for j := 0; j < m; j++ {
                   if s[i][j] == '1' {
                       donors[j] = append(donors[j], i)
                   }
               }
           } else if cnt[i] < per {
               deficit = append(deficit, i)
           }
       }
       // operations: each is [donor, receiver, column]
       ops := make([][3]int, 0)
       // for each deficit row
       for _, i := range deficit {
           for cnt[i] < per {
               found := false
               for j := 0; j < m && !found; j++ {
                   if s[i][j] == '0' {
                       // find donor in this column
                       for idx, d := range donors[j] {
                           if cnt[d] > per {
                               // apply operation: move 1 from d to i at col j
                               cnt[d]--
                               cnt[i]++
                               s[d][j] = '0'
                               s[i][j] = '1'
                               // record op with 1-based indices
                               ops = append(ops, [3]int{d + 1, i + 1, j + 1})
                               // remove donor d from donors[j]
                               donors[j] = append(donors[j][:idx], donors[j][idx+1:]...)
                               found = true
                               break
                           }
                       }
                   }
               }
               // safety: should always find
               if !found {
                   break
               }
           }
       }
       // output
       fmt.Fprintln(writer, len(ops))
       for _, op := range ops {
           fmt.Fprintf(writer, "%d %d %d\n", op[0], op[1], op[2])
       }
   }
}
