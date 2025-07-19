package main

import (
   "bufio"
   "fmt"
   "os"
)

// get returns the minimum number of changes to make the four characters
// form a palindrome under 180-degree rotation of the 2x2 block.
func get(a, b, c, d byte) int {
   // frequency map of the four bytes
   var cnt [256]int
   cnt[a]++
   cnt[b]++
   cnt[c]++
   cnt[d]++
   // collect non-zero frequencies
   freqs := make([]int, 0, 4)
   for _, ch := range []byte{a, b, c, d} {
       if f := cnt[ch]; f > 0 {
           freqs = append(freqs, f)
           cnt[ch] = 0 // avoid duplicates in freqs
       }
   }
   distinct := len(freqs)
   switch distinct {
   case 1:
       return 0
   case 2:
       // two cases: (2,2) => 0, (3,1) or (1,3) => 1
       if freqs[0] == 2 && freqs[1] == 2 {
           return 0
       }
       return 1
   case 3:
       return 1
   case 4:
       return 2
   }
   return 0
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   // read n and two strings
   fmt.Fscan(reader, &n)
   s0 := make([]byte, 0, n)
   s1 := make([]byte, 0, n)
   // scan strings
   var str0, str1 string
   fmt.Fscan(reader, &str0, &str1)
   s0 = []byte(str0)
   s1 = []byte(str1)
   ans := 0
   // process pairs
   for i := 0; i < n/2; i++ {
       ans += get(s0[i], s1[i], s0[n-i-1], s1[n-i-1])
   }
   // middle column if odd
   if n%2 == 1 && s0[n/2] != s1[n/2] {
       ans++
   }
   fmt.Println(ans)
}
