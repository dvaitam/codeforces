package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   // Read heading line
   if !scanner.Scan() {
       return
   }
   s1 := scanner.Text()
   // Read target text line
   if !scanner.Scan() {
       return
   }
   s2 := scanner.Text()

   // Count available letters in the heading (ignore spaces)
   count := make(map[rune]int)
   for _, c := range s1 {
       if c != ' ' {
           count[c]++
       }
   }

   // Try to compose the target text
   ok := true
   for _, c := range s2 {
       if c == ' ' {
           continue
       }
       if count[c] > 0 {
           count[c]--
       } else {
           ok = false
           break
       }
   }

   if ok {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
