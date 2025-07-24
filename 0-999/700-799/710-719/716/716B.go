package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)

   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   if n < 26 {
       fmt.Println(-1)
       return
   }
   arr := []rune(s)
   for i := 0; i+26 <= n; i++ {
       var freq [26]int
       q := 0
       ok := true
       for j := i; j < i+26; j++ {
           if arr[j] == '?' {
               q++
           } else {
               idx := arr[j] - 'A'
               if idx < 0 || idx >= 26 {
                   ok = false
                   break
               }
               if freq[idx] > 0 {
                   ok = false
                   break
               }
               freq[idx]++
           }
       }
       if !ok {
           continue
       }
       used := 0
       for k := 0; k < 26; k++ {
           if freq[k] > 0 {
               used++
           }
       }
       if used+q == 26 {
           // build list of missing letters
           missing := make([]rune, 0, 26-used)
           for k := 0; k < 26; k++ {
               if freq[k] == 0 {
                   missing = append(missing, rune('A'+k))
               }
           }
           // fill in the window
           cur := 0
           for j := i; j < i+26; j++ {
               if arr[j] == '?' {
                   arr[j] = missing[cur]
                   cur++
               }
           }
           // replace remaining '?' with 'A'
           for j := 0; j < n; j++ {
               if arr[j] == '?' {
                   arr[j] = 'A'
               }
           }
           fmt.Println(string(arr))
           return
       }
   }
   fmt.Println(-1)
}
