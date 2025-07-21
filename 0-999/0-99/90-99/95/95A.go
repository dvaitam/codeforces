package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   forbidden := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &forbidden[i])
       forbidden[i] = strings.ToLower(forbidden[i])
   }
   var w string
   fmt.Fscan(reader, &w)
   var letter string
   fmt.Fscan(reader, &letter)
   // target letter in lowercase
   target := letter[0]
   L := len(w)
   covered := make([]bool, L)
   lowerW := strings.ToLower(w)
   // mark covered positions
   for _, pat := range forbidden {
       lp := len(pat)
       for i := 0; i+lp <= L; i++ {
           if lowerW[i:i+lp] == pat {
               for j := i; j < i+lp; j++ {
                   covered[j] = true
               }
           }
       }
   }
   // build result
   out := make([]rune, L)
   for i, ch := range w {
       if covered[i] {
           // replace with target, preserving case
           if ch >= 'A' && ch <= 'Z' {
               // uppercase
               out[i] = rune(strings.ToUpper(letter)[0])
           } else {
               // lowercase or other
               out[i] = rune(target)
           }
       } else {
           out[i] = ch
       }
   }
   fmt.Println(string(out))
}
