package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       // ignore
   }
   // trim newline
   if len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
       // remove all trailing \r or \n
       for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
           s = s[:len(s)-1]
       }
   }
   n := len(s)
   // positions
   pos0 := make([]int, 0, n)
   pos1 := make([]int, 0, n)
   qpos := make([]int, 0, n)
   for i, c := range s {
       switch c {
       case '0':
           pos0 = append(pos0, i)
       case '1':
           pos1 = append(pos1, i)
       case '?':
           qpos = append(qpos, i)
       }
   }
   C0 := len(pos0)
   C1 := len(pos1)
   Q := len(qpos)
   D := n - 2
   if D < 0 {
       D = 0
   }
   m := D / 2
   M := (D + 1) / 2
   // results
   can00 := false
   can11 := false
   can01 := false
   can10 := false
   // check 00 and 11
   if C0 + Q >= m + 2 {
       can00 = true
   }
   if C1 + Q >= M + 2 {
       can11 = true
   }
   // check mixed survivors
   if C0 + Q >= m + 1 && C1 + Q >= M + 1 {
       // L and R bounds for zeros among ?
       L := m + 1 - C0
       if L < 0 {
           L = 0
       }
       R := C1 + Q - (M + 1)
       if R > Q {
           R = Q
       }
       if R >= L {
           // compute Zpos_min
           var Zmin int
           if m < C0 {
               Zmin = pos0[m]
           } else {
               Zmin = qpos[m - C0]
           }
           // compute Opos_max at k=R
           var Omax int
           if M < C1 {
               Omax = pos1[M]
           } else {
               // index in qpos: R + (M - C1)
               idx := R + (M - C1)
               if idx < 0 {
                   idx = 0
               }
               if idx >= len(qpos) {
                   idx = len(qpos) - 1
               }
               Omax = qpos[idx]
           }
           if Zmin < Omax {
               can01 = true
           }
           // compute Zpos_max at k=L
           var Zmax int
           if m < C0 {
               Zmax = pos0[m]
           } else {
               // index: last ? qpos[Q-1]
               Zmax = qpos[Q-1]
           }
           // compute Opos_min
           var Omin int
           if M < C1 {
               Omin = pos1[M]
           } else {
               idx := M - C1
               if idx < 0 {
                   idx = 0
               }
               if idx >= len(qpos) {
                   idx = len(qpos) - 1
               }
               Omin = qpos[idx]
           }
           if Zmax > Omin {
               can10 = true
           }
       }
   }
   // output in lex order
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   if can00 {
       fmt.Fprintln(writer, "00")
   }
   if can01 {
       fmt.Fprintln(writer, "01")
   }
   if can10 {
       fmt.Fprintln(writer, "10")
   }
   if can11 {
       fmt.Fprintln(writer, "11")
   }
}
