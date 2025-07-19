package main

import (
   "bufio"
   "fmt"
   "os"
)

// oddManacher computes array of palindrome radii for odd-length palindromes in pattern
func oddManacher(pattern []byte) []int {
   n := len(pattern)
   radius := make([]int, n)
   loc := 0
   for i := 1; i < n; i++ {
       if i <= loc+radius[loc] {
           mirror := loc - (i - loc)
           // ensure we don't exceed the right boundary
           if radius[mirror] < loc+radius[loc]-i {
               radius[i] = radius[mirror]
           } else {
               radius[i] = loc + radius[loc] - i
           }
       }
       // expand around center i
       for i-radius[i] > 0 && i+radius[i] < n-1 && pattern[i-radius[i]-1] == pattern[i+radius[i]+1] {
           radius[i]++
       }
       if i+radius[i] > loc+radius[loc] {
           loc = i
       }
   }
   return radius
}

// manacher builds the combined Manacher array for string s
func manacher(s []byte) []int {
   n := len(s)
   ext := make([]byte, 2*n+1)
   for i := 0; i < n; i++ {
       ext[2*i+1] = s[i]
   }
   return oddManacher(ext)
}

// isPalindrome checks if s[start:end] is a palindrome using precomputed lengths
func isPalindrome(lengths []int, start, end int) bool {
   // combined index is start+end
   return lengths[start+end] >= end-start
}

func reverseBytes(s []byte) {
   for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
       s[i], s[j] = s[j], s[i]
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var str string
       fmt.Fscan(in, &str)
       s := []byte(str)
       n := len(s)
       l, r := 0, n-1
       for l <= r && s[l] == s[r] {
           l++
           r--
       }
       if r < l {
           // entire string is palindrome
           fmt.Fprintln(out, str)
           continue
       }
       ansLen := 0
       ansStr := ""
       // two passes: original and reversed
       for iter := 0; iter < 2; iter++ {
           lengths := manacher(s)
           best := -1
           for j := l; j <= r; j++ {
               if isPalindrome(lengths, l, j+1) {
                   best = j
               }
           }
           // build candidate
           prefix := string(s[:best+1])
           suffix := string(s[r+1:])
           cand := prefix + suffix
           if len(cand) > ansLen || (len(cand) == ansLen && cand > ansStr) {
               ansLen = len(cand)
               ansStr = cand
           }
           // prepare for reversed pass
           // new l,r on reversed string
           nl := n - 1 - r
           nr := n - 1 - l
           l, r = nl, nr
           reverseBytes(s)
       }
       fmt.Fprintln(out, ansStr)
   }
}
