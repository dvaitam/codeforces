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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   s := make([]byte, 0, 2*n)
   t := make([]byte, 0, 2*n)
   // read s and t
   var ss, tt string
   fmt.Fscan(reader, &ss)
   fmt.Fscan(reader, &tt)
   s = []byte(ss)
   t = []byte(tt)
   // buckets counts
   cnt11, cnt10, cnt01, cnt00 := 0, 0, 0, 0
   for i := 0; i < 2*n; i++ {
       a := s[i] - '0'
       b := t[i] - '0'
       switch a + b {
       case 2:
           cnt11++
       case 1:
           if a == 1 {
               cnt10++
           } else {
               cnt01++
           }
       default:
           cnt00++
       }
   }
   // simulate game
   total := 2 * n
   D := 0
   turn := 0 // 0: first(Y), 1: second(A)
   // process sum2 bucket (1,1)
   for i := 0; i < cnt11; i++ {
       if turn == 0 {
           D += 1
       } else {
           D -= 1
       }
       turn ^= 1
   }
   // process sum1 bucket: type10
   for i := 0; i < cnt10; i++ {
       if turn == 0 {
           D += 1
       }
       // else subtract 0
       turn ^= 1
   }
   // process sum1 bucket: type01
   for i := 0; i < cnt01; i++ {
       if turn == 1 {
           D -= 1
       }
       // else add 0
       turn ^= 1
   }
   // process sum0 bucket
   // remaining turns just flip turn but no score impact
   rem := total - cnt11 - cnt10 - cnt01
   if rem%2 == 1 && turn == 0 {
       // one extra flip if odd, but no score
   }
   // No need to simulate flips for sum0, as no score change

   // determine result
   var res string
   if D > 0 {
       res = "First"
   } else if D < 0 {
       res = "Second"
   } else {
       res = "Draw"
   }
   fmt.Fprintln(writer, res)
}
