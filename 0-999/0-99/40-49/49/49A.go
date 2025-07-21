package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, err := reader.ReadString('\n')
   if err != nil && len(line) == 0 {
       return
   }
   line = strings.TrimRight(line, "\r\n")
   // find last letter before question mark
   var lastLetter rune
   for i := len(line) - 1; i >= 0; i-- {
       ch := rune(line[i])
       if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
           lastLetter = ch
           break
       }
   }
   // normalize to upper
   lastLetter = rune(strings.ToUpper(string(lastLetter))[0])
   // set of vowels
   vowels := "AEIOUY"
   answer := "NO"
   if strings.ContainsRune(vowels, lastLetter) {
       answer = "YES"
   }
   fmt.Println(answer)
}
