package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   pickup  = "PICKUP"
   dropoff = "DROPOFF"
)

func drive(x int) string {
   return fmt.Sprintf("DRIVE %d", x+1)
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func solve(nom, noc, mc, cm, mno, cno []int, writer *bufio.Writer) {
   if len(mc)+len(cm)+len(mno)+len(cno) == 0 {
       fmt.Fprintln(writer, 0)
       writer.Flush()
       os.Exit(0)
   }
   if len(nom) == 0 {
       return
   }
   var ans []string
   // initial drive to first missing
   idx := nom[len(nom)-1]
   nom = nom[:len(nom)-1]
   ans = append(ans, drive(idx))
   mcNext, lastM := true, true
   // alternate pickups
   for {
       if mcNext {
           if len(mc) == 0 {
               break
           }
           ans = append(ans, pickup)
           idx = mc[len(mc)-1]
           mc = mc[:len(mc)-1]
           ans = append(ans, drive(idx))
           ans = append(ans, dropoff)
           lastM = false
           mcNext = false
       } else {
           if len(cm) == 0 {
               break
           }
           ans = append(ans, pickup)
           idx = cm[len(cm)-1]
           cm = cm[:len(cm)-1]
           ans = append(ans, drive(idx))
           ans = append(ans, dropoff)
           lastM = true
           mcNext = true
       }
   }
   // fix one leftover
   if lastM {
       if len(mno) > 0 {
           ans = append(ans, pickup)
           idx = mno[len(mno)-1]
           mno = mno[:len(mno)-1]
           ans = append(ans, drive(idx))
           ans = append(ans, dropoff)
       }
   } else {
       if len(cno) > 0 {
           ans = append(ans, pickup)
           idx = cno[len(cno)-1]
           cno = cno[:len(cno)-1]
           ans = append(ans, drive(idx))
           ans = append(ans, dropoff)
       }
   }
   // pair cross mismatches
   mcc0 := min(len(mc), len(cno))
   cmm0 := min(len(cm), len(mno))
   for i := 0; i < cmm0; i++ {
       if len(noc) == 0 {
           return
       }
       idx = noc[len(noc)-1]
       noc = noc[:len(noc)-1]
       ans = append(ans, drive(idx))
       ans = append(ans, pickup)
       idx = cm[len(cm)-1]
       cm = cm[:len(cm)-1]
       ans = append(ans, drive(idx))
       ans = append(ans, dropoff)
       ans = append(ans, pickup)
       idx = mno[len(mno)-1]
       mno = mno[:len(mno)-1]
       ans = append(ans, drive(idx))
       ans = append(ans, dropoff)
   }
   for i := 0; i < mcc0; i++ {
       if len(nom) == 0 {
           return
       }
       idx = nom[len(nom)-1]
       nom = nom[:len(nom)-1]
       ans = append(ans, drive(idx))
       ans = append(ans, pickup)
       idx = mc[len(mc)-1]
       mc = mc[:len(mc)-1]
       ans = append(ans, drive(idx))
       ans = append(ans, dropoff)
       ans = append(ans, pickup)
       idx = cno[len(cno)-1]
       cno = cno[:len(cno)-1]
       ans = append(ans, drive(idx))
       ans = append(ans, dropoff)
   }
   // move remaining
   for len(mc) > 0 {
       mno = append(mno, mc[len(mc)-1])
       mc = mc[:len(mc)-1]
   }
   for len(cm) > 0 {
       cno = append(cno, cm[len(cm)-1])
       cm = cm[:len(cm)-1]
   }
   // finish leftovers
   for len(mno) > 0 {
       if len(nom) == 0 {
           return
       }
       idx = nom[len(nom)-1]
       nom = nom[:len(nom)-1]
       ans = append(ans, drive(idx))
       ans = append(ans, pickup)
       idx = mno[len(mno)-1]
       mno = mno[:len(mno)-1]
       ans = append(ans, drive(idx))
       ans = append(ans, dropoff)
   }
   for len(cno) > 0 {
       if len(noc) == 0 {
           return
       }
       idx = noc[len(noc)-1]
       noc = noc[:len(noc)-1]
       ans = append(ans, drive(idx))
       ans = append(ans, pickup)
       idx = cno[len(cno)-1]
       cno = cno[:len(cno)-1]
       ans = append(ans, drive(idx))
       ans = append(ans, dropoff)
   }
   // output answer
   fmt.Fprintln(writer, len(ans))
   for _, op := range ans {
       fmt.Fprintln(writer, op)
   }
   writer.Flush()
   os.Exit(0)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s, t string
   fmt.Fscan(reader, &s, &t)
   var nom, noc, mc, cm, mno, cno []int
   for i := 0; i < n; i++ {
       if s[i] == t[i] {
           continue
       }
       switch s[i] {
       case '-':
           if t[i] == 'M' {
               nom = append(nom, i)
           } else if t[i] == 'C' {
               noc = append(noc, i)
           }
       case 'C':
           if t[i] == 'M' {
               cm = append(cm, i)
           } else if t[i] == '-' {
               cno = append(cno, i)
           }
       case 'M':
           if t[i] == 'C' {
               mc = append(mc, i)
           } else if t[i] == '-' {
               mno = append(mno, i)
           }
       }
   }
   solve(nom, noc, mc, cm, mno, cno, writer)
   solve(noc, nom, cm, mc, cno, mno, writer)
}
