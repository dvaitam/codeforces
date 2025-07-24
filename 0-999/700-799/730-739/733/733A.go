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
   // Vowel letters including 'Y'
   vowels := map[byte]bool{
       'A': true, 'E': true, 'I': true,
       'O': true, 'U': true, 'Y': true,
   }
   // positions including start (0) and end (len(s)+1)
   var pos []int
   pos = append(pos, 0)
   for i := 0; i < len(s); i++ {
       if vowels[s[i]] {
           pos = append(pos, i+1)
       }
   }
   pos = append(pos, len(s)+1)
   // compute maximum gap
   ans := 0
   for i := 1; i < len(pos); i++ {
       gap := pos[i] - pos[i-1]
       if gap > ans {
           ans = gap
       }
   }
   fmt.Println(ans)
}
