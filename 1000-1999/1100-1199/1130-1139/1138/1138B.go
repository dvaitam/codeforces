package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   _, err := fmt.Fscan(reader, &n)
   if err != nil {
       return
   }
   // Read strings a and b
   var aStr, bStr string
   // Consume next tokens as strings (could be without spaces)
   fmt.Fscan(reader, &aStr)
   fmt.Fscan(reader, &bStr)
   a := make([]int, n)
   b := make([]int, n)
   v := make([]int, n)
   a0, a1, a2 := 0, 0, 0
   sumb := 0
   for i := 0; i < n; i++ {
       if i < len(aStr) {
           a[i] = int(aStr[i] - '0')
       }
       if i < len(bStr) {
           b[i] = int(bStr[i] - '0')
       }
       if b[i] == 1 {
           sumb++
       }
       v[i] = a[i] + b[i]
       switch v[i] {
       case 0:
           a0++
       case 1:
           a1++
       case 2:
           a2++
       }
   }
   half := n / 2
   // Try choose cnt2=i, cnt1=j
   for i2 := 0; i2 <= a2; i2++ {
       for j1 := 0; j1 <= a1; j1++ {
           if i2+j1 > half {
               continue
           }
           if 2*i2+j1 == sumb {
               cnt2, cnt1 := i2, j1
               cnt0 := half - i2 - j1
               if cnt0 > a0 {
                   continue
               }
               out := bufio.NewWriter(os.Stdout)
               defer out.Flush()
               for k := 0; k < n; k++ {
                   if v[k] == 1 && cnt1 > 0 {
                       fmt.Fprintf(out, "%d ", k+1)
                       cnt1--
                   } else if v[k] == 2 && cnt2 > 0 {
                       fmt.Fprintf(out, "%d ", k+1)
                       cnt2--
                   } else if v[k] == 0 && cnt0 > 0 {
                       fmt.Fprintf(out, "%d ", k+1)
                       cnt0--
                   }
               }
               return
           }
       }
   }
   // No solution
   fmt.Print("-1")
}
