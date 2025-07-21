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

   var a, b string
   fmt.Fscan(reader, &a)
   fmt.Fscan(reader, &b)
   n, m := len(a), len(b)
   if a == b {
       fmt.Fprintln(writer, "YES")
       return
   }
   // base for rolling hash
   const B uint64 = 91138233
   // precompute powers
   pow := make([]uint64, m)
   if m > 0 {
       pow[0] = 1
       for i := 1; i < m; i++ {
           pow[i] = pow[i-1] * B
       }
   }
   // hash of b
   var hb uint64
   for i := 0; i < m; i++ {
       if b[i] == '1' {
           hb += pow[m-1-i]
       }
   }
   // try each possible starting point after removals
   for k := 0; k < n; k++ {
       // build initial window s of length m
       s := []byte{}
       // copy suffix
       if n-k >= m {
           s = []byte(a[k : k+m])
       } else {
           // shorter: copy all and append parity bits
           s = []byte(a[k:])
           // compute initial parity
           var parity byte
           for _, c := range s {
               parity ^= c - '0'
           }
           // extend until length m
           for len(s) < m {
               s = append(s, '0'+parity)
               // after appending parity bit equal to old parity, new parity is 0
               parity = 0
           }
       }
       // compute hash and parity of s
       var hs uint64
       var par byte
       for i := 0; i < m; i++ {
           if s[i] == '1' {
               hs += pow[m-1-i]
               par ^= 1
           }
       }
       // sliding window simulation (remove front, append parity)
       for t := 0; t <= m; t++ {
           if hs == hb {
               // verify to avoid rare collision
               if string(s) == b {
                   fmt.Fprintln(writer, "YES")
                   return
               }
           }
           if t == m {
               break
           }
           // parity to append is current parity
           p := par
           // remove front
           front := s[0] - '0'
           // update rolling hash: drop front bit then shift and add p
           hs = (hs - uint64(front)*pow[m-1]) * B
           if p == 1 {
               hs += 1
           }
           // slide bytes
           copy(s[0:m-1], s[1:m])
           s[m-1] = '0' + p
           // update parity: new parity is old front bit
           par = front
       }
   }
   fmt.Fprintln(writer, "NO")
}
