package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read target string s
   sLine, err := reader.ReadString('\n')
   if err != nil && len(sLine) == 0 {
       return
   }
   s := strings.TrimSpace(sLine)
   // Read source string t
   tLine, err := reader.ReadString('\n')
   if err != nil && len(tLine) == 0 {
       return
   }
   t := strings.TrimSpace(tLine)

   // Count available letters in t by case
   var cntLower [26]int
   var cntUpper [26]int
   for _, ch := range t {
       if ch >= 'a' && ch <= 'z' {
           cntLower[ch-'a']++
       } else if ch >= 'A' && ch <= 'Z' {
           cntUpper[ch-'A']++
       }
   }

   yay := 0
   whoops := 0
   // Build message for s
   for _, ch := range s {
       if ch >= 'a' && ch <= 'z' {
           idx := ch - 'a'
           if cntLower[idx] > 0 {
               yay++
               cntLower[idx]--
           } else if cntUpper[idx] > 0 {
               whoops++
               cntUpper[idx]--
           }
       } else if ch >= 'A' && ch <= 'Z' {
           idx := ch - 'A'
           if cntUpper[idx] > 0 {
               yay++
               cntUpper[idx]--
           } else if cntLower[idx] > 0 {
               whoops++
               cntLower[idx]--
           }
       }
   }

   fmt.Printf("%d %d", yay, whoops)
}
