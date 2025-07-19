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
   if n < 25 {
       fmt.Println(-1)
       return
   }
   y := 0
   for i := 5; i*i <= n; i++ {
       if n%i == 0 {
           y = i
           break
       }
   }
   if y == 0 {
       fmt.Println(-1)
       return
   }
   other := n / y
   if other < 5 {
       fmt.Println(-1)
       return
   }
   if y == 5 {
       patterns := []string{"aeiou", "eioua", "iouae", "ouaei", "uaeio"}
       var builder strings.Builder
       for i := 0; i < other; i++ {
           builder.WriteString(patterns[i%5])
       }
       fmt.Print(builder.String())
       return
   }
   vowels := []byte{'a', 'e', 'i', 'o', 'u'}
   var builder strings.Builder
   x := 0
   for j := 0; j < other; j++ {
       for i := 0; i < y; i++ {
           builder.WriteByte(vowels[x%5])
           x++
       }
   }
   fmt.Print(builder.String())
}
