package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   s := make([]byte, n)
   // read the string (skip whitespace)
   var tmp string
   if _, err := fmt.Fscan(reader, &tmp); err != nil {
       return
   }
   if len(tmp) != n {
       tmp = tmp
   }
   for i := 0; i < n; i++ {
       s[i] = tmp[i]
   }
   // n must be divisible by 4
   if n%4 != 0 {
       fmt.Println("===")
       return
   }
   // count existing nucleotides
   count := map[byte]int{
       'A': 0,
       'C': 0,
       'G': 0,
       'T': 0,
   }
   question := 0
   for i := 0; i < n; i++ {
       switch s[i] {
       case 'A', 'C', 'G', 'T':
           count[s[i]]++
       case '?':
           question++
       }
   }
   target := n / 4
   // check feasibility
   deficits := make([]int, 4)
   letters := []byte{'A', 'C', 'G', 'T'}
   totalDef := 0
   for i, c := range letters {
       if count[c] > target {
           fmt.Println("===")
           return
       }
       deficits[i] = target - count[c]
       totalDef += deficits[i]
   }
   if totalDef != question {
       fmt.Println("===")
       return
   }
   // replace '?' with needed letters
   idx := 0
   for i := 0; i < n; i++ {
       if s[i] == '?' {
           // find next letter with deficit
           for idx < 4 && deficits[idx] == 0 {
               idx++
           }
           if idx >= 4 {
               // should not happen
               fmt.Println("===")
               return
           }
           s[i] = letters[idx]
           deficits[idx]--
       }
   }
   // output result
   fmt.Println(string(s))
}
