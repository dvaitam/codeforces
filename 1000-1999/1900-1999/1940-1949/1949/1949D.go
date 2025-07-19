package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

type pair struct{ i, j int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   str := make([][]rune, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       str[i] = []rune(s)
   }
   var e [2][]pair
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           switch str[i][j] {
           case 'F':
               if i < j {
                   e[0] = append(e[0], pair{i, j})
               }
           case 'S':
               if i < j {
                   e[1] = append(e[1], pair{i, j})
               }
           }
       }
   }
   full := (1 << n) - 1
   lim := (3*n + 3) / 4
   symbols := []rune{'F', 'S'}
   // iterate all masks
   for mask := 0; mask < 1<<n; mask++ {
       cnt1 := bits.OnesCount(uint(mask))
       if cnt1-1 > lim {
           continue
       }
       zeros := bits.OnesCount(uint(mask ^ full))
       if zeros*2 > lim {
           continue
       }
       for x := 0; x < 2; x++ {
           bad := false
           for _, p := range e[x] {
               a := (mask >> p.i) & 1
               b := (mask >> p.j) & 1
               if a^b != 0 {
                   bad = true
                   break
               }
           }
           if bad {
               continue
           }
           c0 := 0
           for _, p := range e[1-x] {
               if ((mask>>p.i)&1) == 1 && ((mask>>p.j)&1) == 1 {
                   c0++
               }
           }
           if zeros*2 + c0 > lim {
               continue
           }
           // build result
           for i := 0; i < n; i++ {
               for j := i + 1; j < n; j++ {
                   if str[i][j] == '?' {
                       a := (mask >> i) & 1
                       b := (mask >> j) & 1
                       var ch rune
                       if a == 1 && b == 1 {
                           ch = symbols[x]
                       } else if a^b == 1 {
                           ch = symbols[1-x]
                       } else {
                           ch = 'F'
                       }
                       str[i][j], str[j][i] = ch, ch
                   }
               }
           }
           // output
           for i := 0; i < n; i++ {
               writer.WriteString(string(str[i]))
               writer.WriteByte('\n')
           }
           return
       }
   }
}
