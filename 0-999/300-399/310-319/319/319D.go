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

   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   b := []byte(s)
   // Greedy: repeatedly remove the shortest repeating block (XX) at leftmost
   for {
       n := len(b)
       found := false
       // L is length of X
       for L := 1; 2*L <= n; L++ {
           limit := n - 2*L + 1
           for i := 0; i < limit; i++ {
               // compare b[i:i+L] and b[i+L:i+2L]
               eq := true
               for k := 0; k < L; k++ {
                   if b[i+k] != b[i+L+k] {
                       eq = false
                       break
                   }
               }
               if eq {
                   // remove the second X (i.e., b[i+L:i+2L])
                   // keep b[0:i+L] + b[i+2L:]
                   nb := make([]byte, 0, n- L)
                   nb = append(nb, b[:i+L]...)
                   nb = append(nb, b[i+2*L:]...)
                   b = nb
                   found = true
                   break
               }
           }
           if found {
               break
           }
       }
       if !found {
           break
       }
   }
   writer.Write(b)
}
