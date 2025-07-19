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

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for t := 0; t < T; t++ {
       var n int
       fmt.Fscan(reader, &n)
       var s1, s2 string
       fmt.Fscan(reader, &s1)
       fmt.Fscan(reader, &s2)
       // detect if both matching and mismatching positions exist
       fl1, fl2 := false, false
       for i := 0; i < n; i++ {
           if s1[i] == s2[i] {
               fl1 = true
           } else {
               fl2 = true
           }
       }
       if fl1 && fl2 {
           fmt.Fprintln(writer, "NO")
           continue
       }
       // build operations
       type op struct{ l, r int }
       var ops []op
       k := 0
       for i := 0; i < n; i++ {
           if s1[i] == '1' {
               ops = append(ops, op{i + 1, i + 1})
               if i >= 1 {
                   k++
               }
           }
       }
       // check parity
       if ((int(s2[0]-'0') + k) & 1) != 0 {
           // add three specific operations
           ops = append(ops, op{1, n - 1})
           ops = append(ops, op{n, n})
           ops = append(ops, op{1, n})
       }
       // output
       fmt.Fprintln(writer, "YES")
       fmt.Fprintln(writer, len(ops))
       for _, e := range ops {
           fmt.Fprintln(writer, e.l, e.r)
       }
   }
}
