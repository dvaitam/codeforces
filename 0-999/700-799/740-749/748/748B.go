package main

import (
   "fmt"
)

func main() {
   var s, t string
   if _, err := fmt.Scan(&s, &t); err != nil {
       return
   }
   if len(s) != len(t) {
       fmt.Println(-1)
       return
   }
   const alpha = 26
   f := make([]int, alpha)
   for i := 0; i < alpha; i++ {
       f[i] = -1
   }
   for i := 0; i < len(s); i++ {
       a := int(s[i] - 'a')
       b := int(t[i] - 'a')
       if f[a] == -1 && f[b] == -1 {
           if a != b {
               f[a] = b
               f[b] = a
           } else {
               f[a] = a
           }
       } else if f[a] != -1 {
           if f[a] != b {
               fmt.Println(-1)
               return
           }
           if f[b] == -1 {
               f[b] = a
           } else if f[b] != a {
               fmt.Println(-1)
               return
           }
       } else {
           // f[a] == -1 && f[b] != -1
           if f[b] != a {
               fmt.Println(-1)
               return
           }
           f[a] = b
       }
   }
   // collect swap pairs
   type pair struct{ x, y byte }
   var pairs []pair
   for c := 0; c < alpha; c++ {
       if f[c] == -1 {
           f[c] = c
       }
   }
   for c := 0; c < alpha; c++ {
       d := f[c]
       if d > c {
           pairs = append(pairs, pair{byte(c), byte(d)})
       }
   }
   fmt.Println(len(pairs))
   for _, p := range pairs {
       fmt.Printf("%c %c\n", 'a'+p.x, 'a'+p.y)
   }
}
