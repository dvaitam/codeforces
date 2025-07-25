package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // total intervals
   total := int64(n) * int64(n+1) / 2
   // count small intervals of length up to 8
   var smallCount int64
   maxL := 8
   if n < maxL {
       maxL = n
   }
   for L := 1; L <= maxL; L++ {
       smallCount += int64(n - L + 1)
   }
   // intervals of length >=9 automatically contain a monochromatic 3-term AP (van der Waerden W(2,3)=9)
   longCount := total - smallCount
   // count small intervals that actually contain a progression
   sBytes := []byte(s)
   var smallGood int64
   for L := 1; L <= maxL; L++ {
       for l := 0; l+L <= n; l++ {
           r := l + L - 1
           found := false
           // try all k where 2*k <= L-1
           for k := 1; 2*k <= L-1; k++ {
               for x := l; x+2*k <= r; x++ {
                   b := sBytes[x]
                   if b == sBytes[x+k] && b == sBytes[x+2*k] {
                       found = true
                       break
                   }
               }
               if found {
                   break
               }
           }
           if found {
               smallGood++
           }
       }
   }
   ans := longCount + smallGood
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprint(w, ans)
}
