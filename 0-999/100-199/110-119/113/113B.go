package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s, b, e string
   if _, err := fmt.Fscan(reader, &s, &b, &e); err != nil {
       return
   }
   ns, nb, ne := len(s), len(b), len(e)
   nn := nb
   if ne > nn {
       nn = ne
   }

   fb := make([]bool, ns)
   fe := make([]bool, ns)
   // mark prefix and suffix match positions
   for i := 0; i+nb <= ns; i++ {
       if s[i:i+nb] == b {
           fb[i] = true
       }
   }
   for i := 0; i+ne <= ns; i++ {
       if s[i:i+ne] == e {
           fe[i] = true
       }
   }

   // compute longest common prefix lengths for all suffix pairs
   p := make([]int, ns+1)
   q := make([]int, ns+1)
   h := make([]int, ns)
   for i := 0; i < ns; i++ {
       for j := i - 1; j >= 0; j-- {
           if s[i] == s[j] {
               q[j+1] = p[j] + 1
           } else {
               q[j+1] = 0
           }
           if q[j+1] > h[i] {
               h[i] = q[j+1]
           }
       }
       // swap p and q for next iteration
       p, q = q, p
   }

   // count valid substrings
   w := 0
   for i := ns - 1; i >= 0; i-- {
       if !fb[i] {
           continue
       }
       for j := i + nn; j <= ns; j++ {
           if fe[j-ne] && h[j-1] < j-i {
               w++
           }
       }
   }
   fmt.Println(w)
