package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && len(s) == 0 {
       return
   }
   seen := make(map[rune]bool)
   for _, c := range s {
       if c >= 'a' && c <= 'z' {
           seen[c] = true
       }
   }
   fmt.Println(len(seen))
}
