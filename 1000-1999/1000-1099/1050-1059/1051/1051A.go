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
   fmt.Fscan(reader, &T)
   for t := 0; t < T; t++ {
       var s string
       fmt.Fscan(reader, &s)
       b := []byte(s)
       var U, D, L []int
       for i, c := range b {
           switch {
           case c >= 'A' && c <= 'Z':
               U = append(U, i)
           case c >= '0' && c <= '9':
               D = append(D, i)
           default:
               L = append(L, i)
           }
       }
       mU := len(U) == 0
       mD := len(D) == 0
       mL := len(L) == 0
       numMissing := 0
       if mU {
           numMissing++
       }
       if mD {
           numMissing++
       }
       if mL {
           numMissing++
       }
       switch numMissing {
       case 1:
           var rep byte
           if mU {
               rep = 'A'
           } else if mD {
               rep = '4'
           } else {
               rep = 'a'
           }
           if len(U) > 1 {
               b[U[0]] = rep
           } else if len(D) > 1 {
               b[D[0]] = rep
           } else if len(L) > 1 {
               b[L[0]] = rep
           }
       case 2:
           var missing []byte
           if mU {
               missing = append(missing, 'A')
           }
           if mD {
               missing = append(missing, '4')
           }
           if mL {
               missing = append(missing, 'a')
           }
           if len(U) > 0 {
               b[U[0]] = missing[0]
               b[U[1]] = missing[1]
           } else if len(D) > 0 {
               b[D[0]] = missing[0]
               b[D[1]] = missing[1]
           } else if len(L) > 0 {
               b[L[0]] = missing[0]
               b[L[1]] = missing[1]
           }
       }
       fmt.Fprintln(writer, string(b))
   }
}
