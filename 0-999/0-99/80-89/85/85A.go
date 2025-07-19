package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   dat := make([][]byte, 4)
   for i := 0; i < 4; i++ {
       dat[i] = make([]byte, n)
   }

   if n%2 == 0 {
       t := 0
       for i := 0; i < n-1; i += 2 {
           ch := byte('a' + t)
           dat[0][i], dat[0][i+1] = ch, ch
           chg := byte('g' + t)
           dat[3][i], dat[3][i+1] = chg, chg
           t = (t + 1) % 2
       }
       t = 0
       for i := 1; i < n-2; i += 2 {
           ch := byte('c' + t)
           dat[1][i], dat[1][i+1] = ch, ch
           che := byte('e' + t)
           dat[2][i], dat[2][i+1] = che, che
           t = (t + 1) % 2
       }
       if n > 0 {
           dat[1][0], dat[2][0] = 'y', 'y'
           dat[1][n-1], dat[2][n-1] = 'z', 'z'
       }
   } else {
       t := 0
       for i := 0; i < n-1; i += 2 {
           ch := byte('a' + t)
           dat[0][i], dat[0][i+1] = ch, ch
           chg := byte('g' + t)
           dat[1][i], dat[1][i+1] = chg, chg
           t = (t + 1) % 2
       }
       t = 0
       for i := 1; i <= n-2; i += 2 {
           ch := byte('c' + t)
           dat[2][i], dat[2][i+1] = ch, ch
           che := byte('e' + t)
           dat[3][i], dat[3][i+1] = che, che
           t = (t + 1) % 2
       }
       if n > 0 {
           dat[2][0], dat[3][0] = 'y', 'y'
           dat[0][n-1], dat[1][n-1] = 'z', 'z'
       }
   }

   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < 4; i++ {
       writer.Write(dat[i])
       writer.WriteByte('\n')
   }
}
