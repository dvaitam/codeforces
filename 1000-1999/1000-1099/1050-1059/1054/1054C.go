package main

import (
   "bufio"
   "fmt"
   "os"
)

type node struct {
   w, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   L := make([]int, n)
   R := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &L[i])
       if L[i] > i {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &R[i])
       if R[i] > n-1-i {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   seq := make([]int, n)
   cnt := n
   var Q []node
   // initial positions with zero constraints
   for i := 0; i < n; i++ {
       if L[i] == 0 && R[i] == 0 {
           seq[i] = cnt
           Q = append(Q, node{cnt, i})
       }
   }
   // assign values
   for len(Q) > 0 {
       // process current layer
       for len(Q) > 0 {
           x := Q[0]
           Q = Q[1:]
           for j := x.id + 1; j < n; j++ {
               if seq[j] != 0 {
                   continue
               }
               L[j]--
           }
           for j := 0; j < x.id; j++ {
               if seq[j] != 0 {
                   continue
               }
               R[j]--
           }
       }
       cnt--
       // find new zeros
       for i := 0; i < n; i++ {
           if seq[i] == 0 && L[i] == 0 && R[i] == 0 {
               seq[i] = cnt
               Q = append(Q, node{cnt, i})
           }
       }
   }
   // check
   for i := 0; i < n; i++ {
       if seq[i] == 0 {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   // output
   fmt.Fprintln(writer, "YES")
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, seq[i])
   }
   writer.WriteByte('\n')
}
