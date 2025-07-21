package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var k int
   var s string
   if _, err := fmt.Fscan(in, &k); err != nil {
       return
   }
   if _, err := fmt.Fscan(in, &s); err != nil {
       return
   }
   n := len(s)
   res := []rune(s)
   // initial mirror
   for i := 0; i < n/2; i++ {
       j := n - 1 - i
       if res[i] != '?' && res[j] != '?' {
           if res[i] != res[j] {
               fmt.Println("IMPOSSIBLE")
               return
           }
       } else if res[i] == '?' && res[j] != '?' {
           res[i] = res[j]
       } else if res[i] != '?' && res[j] == '?' {
           res[j] = res[i]
       }
   }
   // track used letters
   used := make([]bool, k)
   for _, c := range res {
       if c >= 'a' && int(c-'a') < k {
           used[c-'a'] = true
       }
   }
   // collect empty slots (pairs)
   type slot struct{ i, j int }
   var slots []slot
   for i := 0; i < n/2; i++ {
       j := n - 1 - i
       if res[i] == '?' && res[j] == '?' {
           slots = append(slots, slot{i, j})
       }
   }
   // center
   hasCenter := n%2 == 1 && res[n/2] == '?'
   if hasCenter {
       slots = append(slots, slot{n / 2, n / 2})
   }
   // missing letters
   var missing []rune
   for i := 0; i < k; i++ {
       if !used[i] {
           missing = append(missing, rune('a'+i))
       }
   }
   if len(missing) > len(slots) {
       fmt.Println("IMPOSSIBLE")
       return
   }
   // fill missing letters: lexicographically minimal => assign smallest letters to earliest slots
   for idx, c := range missing {
       sl := slots[idx]
       res[sl.i] = c
       res[sl.j] = c
   }
   // fill remaining with 'a'
   for idx := len(missing); idx < len(slots); idx++ {
       sl := slots[idx]
       res[sl.i] = 'a'
       res[sl.j] = 'a'
   }
   // final check: no '?'
   for _, c := range res {
       if c == '?' {
           fmt.Println("IMPOSSIBLE")
           return
       }
   }
   // ensure all used
   finalUsed := make([]bool, k)
   for _, c := range res {
       if c >= 'a' && int(c-'a') < k {
           finalUsed[c-'a'] = true
       }
   }
   for i := 0; i < k; i++ {
       if !finalUsed[i] {
           fmt.Println("IMPOSSIBLE")
           return
       }
   }
   fmt.Println(string(res))
}
