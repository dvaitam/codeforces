package main

import (
   "bufio"
   "os"
   "strconv"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   b, _ := reader.ReadByte()
   for b < '0' || b > '9' {
       b, _ = reader.ReadByte()
   }
   x := 0
   for b >= '0' && b <= '9' {
       x = x*10 + int(b-'0')
       b, _ = reader.ReadByte()
   }
   return x
}

func main() {
   defer writer.Flush()
   t := readInt()
   for ; t > 0; t-- {
       n := readInt()
       m := readInt()
       used := make([]bool, 3*n+1)
       matching := make([]int, 0, n)
       for i := 1; i <= m; i++ {
           u := readInt()
           v := readInt()
           if len(matching) < n && !used[u] && !used[v] {
               used[u] = true
               used[v] = true
               matching = append(matching, i)
           }
       }
       if len(matching) >= n {
           writer.WriteString("Matching\n")
           for i := 0; i < n; i++ {
               writer.WriteString(strconv.Itoa(matching[i]))
               if i+1 < n {
                   writer.WriteByte(' ')
               }
           }
           writer.WriteByte('\n')
       } else {
           ind := make([]int, 0, n)
           for v := 1; v <= 3*n && len(ind) < n; v++ {
               if !used[v] {
                   ind = append(ind, v)
               }
           }
           if len(ind) >= n {
               writer.WriteString("IndSet\n")
               for i := 0; i < n; i++ {
                   writer.WriteString(strconv.Itoa(ind[i]))
                   if i+1 < n {
                       writer.WriteByte(' ')
                   }
               }
               writer.WriteByte('\n')
           } else {
               writer.WriteString("Impossible\n")
           }
       }
   }
}
