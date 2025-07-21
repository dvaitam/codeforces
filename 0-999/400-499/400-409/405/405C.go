package main

import (
   "io"
   "os"
)

func main() {
   data, err := io.ReadAll(os.Stdin)
   if err != nil {
       os.Exit(1)
   }
   pos := 0
   n, p := nextInt(data, pos)
   pos = p
   cur := 0
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           x, p2 := nextInt(data, pos)
           pos = p2
           if i == j && x == 1 {
               cur ^= 1
           }
       }
   }
   q, p3 := nextInt(data, pos)
   pos = p3
   res := make([]byte, 0, q)
   for k := 0; k < q; k++ {
       t, p4 := nextInt(data, pos)
       pos = p4
       if t == 3 {
           res = append(res, byte('0'+cur))
       } else {
           // skip index for row/column flip
           _, p5 := nextInt(data, pos)
           pos = p5
           cur ^= 1
       }
   }
   os.Stdout.Write(res)
}

// nextInt parses next non-negative integer from data starting at idx
func nextInt(data []byte, idx int) (int, int) {
   n := len(data)
   // skip non-digits
   for idx < n && (data[idx] < '0' || data[idx] > '9') {
       idx++
   }
   x := 0
   for idx < n && data[idx] >= '0' && data[idx] <= '9' {
       x = x*10 + int(data[idx]-'0')
       idx++
   }
   return x, idx
}
