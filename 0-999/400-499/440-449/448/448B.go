package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s, t string
   // read s
   if _, err := fmt.Fscanln(reader, &s); err != nil {
       return
   }
   // read t
   if _, err := fmt.Fscanln(reader, &t); err != nil {
       return
   }
   // frequency counts
   freqS := [26]int{}
   freqT := [26]int{}
   for _, ch := range s {
       freqS[ch-'a']++
   }
   for _, ch := range t {
       freqT[ch-'a']++
   }
   // check if s has enough letters to form t
   for i := 0; i < 26; i++ {
       if freqS[i] < freqT[i] {
           fmt.Println("need tree")
           return
       }
   }
   // check if t is subsequence of s
   j := 0
   for i := 0; i < len(s) && j < len(t); i++ {
       if s[i] == t[j] {
           j++
       }
   }
   if j == len(t) {
       // t is subsequence of s
       if len(s) != len(t) {
           fmt.Println("automaton")
       } else {
           // equal length and subsequence implies equal strings, but s != t per constraints
           // fall back to array
           fmt.Println("array")
       }
       return
   }
   // not a subsequence
   if len(s) == len(t) {
       fmt.Println("array")
   } else {
       fmt.Println("both")
   }
}
