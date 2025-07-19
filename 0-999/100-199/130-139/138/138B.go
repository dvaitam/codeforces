package main

import (
   "bufio"
   "fmt"
   "os"
)

func reverse(b []byte) {
   for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
       b[i], b[j] = b[j], b[i]
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   cnt := [10]int{}
   for i := 0; i < len(s); i++ {
       cnt[s[i]-'0']++
   }
   a1 := make([]byte, 0, len(s))
   a2 := make([]byte, 0, len(s))
   // pre-pair zeros while cnt[0] > cnt[9]
   for cnt[0] > cnt[9] {
       cnt[0]--
       a1 = append(a1, '0')
       a2 = append(a2, '0')
   }
   // helper to compute number of carries for initial pair (i,j)
   find := func(i, j int) int {
       c1 := cnt
       c2 := cnt
       if c1[i] == 0 || c2[j] == 0 {
           return 0
       }
       c1[i]--
       c2[j]--
       ret := 1
       for x := 0; x < 10; x++ {
           for y := 0; y < 10; y++ {
               if x+y == 9 {
                   if c1[x] < c2[y] {
                       ret += c1[x]
                   } else {
                       ret += c2[y]
                   }
               }
           }
       }
       return ret
   }
   best := 0
   mi, mj := -1, -1
   for i := 0; i < 10; i++ {
       for j := 0; j < 10; j++ {
           if i+j != 10 {
               continue
           }
           if x := find(i, j); x > best {
               best = x
               mi, mj = i, j
           }
       }
   }
   // rebuild
   if mi != -1 && mj != -1 {
       a1 = append(a1, byte('0'+mi))
       a2 = append(a2, byte('0'+mj))
       c1 := cnt
       c2 := cnt
       c1[mi]--
       c2[mj]--
       // pair digits summing to 9
       for i := 0; i < 10; i++ {
           for j := 0; j < 10; j++ {
               if i+j == 9 {
                   for c1[i] > 0 && c2[j] > 0 {
                       c1[i]--
                       c2[j]--
                       a1 = append(a1, byte('0'+i))
                       a2 = append(a2, byte('0'+j))
                   }
               }
           }
       }
       // append remaining
       for d := 0; d < 10; d++ {
           for c1[d] > 0 {
               c1[d]--
               a1 = append(a1, byte('0'+d))
           }
       }
       for d := 0; d < 10; d++ {
           for c2[d] > 0 {
               c2[d]--
               a2 = append(a2, byte('0'+d))
           }
       }
       reverse(a1)
       reverse(a2)
       fmt.Println(string(a1))
       fmt.Println(string(a2))
   } else {
       if len(a1) == 0 {
           // no pairing done
           fmt.Println(s)
           fmt.Println(s)
       } else {
           // append remaining and print
           for d := 0; d < 10; d++ {
               for cnt[d] > 0 {
                   cnt[d]--
                   a1 = append(a1, byte('0'+d))
               }
           }
           for d := 0; d < 10; d++ {
               for cnt[d] > 0 {
                   cnt[d]--
                   a2 = append(a2, byte('0'+d))
               }
           }
           reverse(a1)
           reverse(a2)
           fmt.Println(string(a1))
           fmt.Println(string(a2))
       }
   }
}
