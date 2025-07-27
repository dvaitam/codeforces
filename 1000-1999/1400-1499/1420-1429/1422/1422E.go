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
   n := len(s)
   // next position where char changes
   nextDiff := make([]int, n)
   for i := n - 1; i >= 0; i-- {
       if i+1 < n && s[i] == s[i+1] {
           nextDiff[i] = nextDiff[i+1]
       } else {
           nextDiff[i] = i + 1
       }
   }
   // DP arrays for f[i]
   isEmpty := make([]bool, n+1)
   headChar := make([]byte, n+1)
   firstDiff := make([]byte, n+1)
   headRunLen := make([]int, n+1)
   length := make([]int, n+1)
   pref := make([]string, n+1)
   suff := make([]string, n+1)
   // base f[n] is empty
   isEmpty[n] = true
   length[n] = 0
   // DP
   for i := n - 1; i >= 0; i-- {
       c := s[i]
       k := nextDiff[i] - i
       r := i + k
       // chooseSmall: whether keep minimal t
       var chooseSmall bool
       if r >= n || isEmpty[r] {
           chooseSmall = true
       } else {
           hc := headChar[r]
           if hc != c {
               chooseSmall = hc < c
           } else {
               fd := firstDiff[r]
               if fd == 0 {
                   // all equal
                   chooseSmall = true
               } else {
                   chooseSmall = fd < c
               }
           }
       }
       // minimal and maximal t
       tMin := k & 1
       tMax := k
       t := tMax
       if chooseSmall {
           t = tMin
       }
       // compute f[i]
       length[i] = t
       if r <= n {
           length[i] += length[r]
       }
       isEmpty[i] = (length[i] == 0)
       // prefix up to 5
       // build pref[i]
       maxPref := 5
       if length[i] <= maxPref {
           // full prefix
           // build full string up to length[i]
           var tmp []byte
           for j := 0; j < t; j++ {
               tmp = append(tmp, c)
           }
           if !isEmpty[r] {
               // append up to remaining
               want := length[i] - t
               curPref := pref[r]
               if len(curPref) > want {
                   curPref = curPref[:want]
               }
               tmp = append(tmp, curPref...)
           }
           pref[i] = string(tmp)
       } else {
           // need first 5
           if t >= maxPref {
               pref[i] = string(make([]byte, maxPref, maxPref)) // fill later
               // fill with c
               buf := []byte(pref[i])
               for j := range buf {
                   buf[j] = c
               }
               pref[i] = string(buf)
           } else {
               // t < 5
               var buf []byte
               for j := 0; j < t; j++ {
                   buf = append(buf, c)
               }
               need := maxPref - t
               // take from pref[r]
               add := pref[r]
               if len(add) > need {
                   add = add[:need]
               }
               buf = append(buf, add...)
               pref[i] = string(buf)
           }
       }
       // suffix up to 2
       // build suff[i]
       if length[i] <= 2 {
           // full suffix
           var tmp []byte
           // total full = c^t + full r
           for j := 0; j < t; j++ {
               tmp = append(tmp, c)
           }
           if !isEmpty[r] {
               // take full from suff[r], but suff[r] may contain just full
               add := suff[r]
               tmp = append(tmp, add...)
           }
           suff[i] = string(tmp)
       } else {
           // length > 2
           if r < n && length[r] >= 2 {
               suff[i] = suff[r]
           } else if r < n && length[r] == 1 {
               // one from R and one from c
               // order: c then R[0]
               tmp := []byte{c, suff[r][0]}
               suff[i] = string(tmp)
           } else {
               // r>=n or length[r]==0: take two c's
               suff[i] = string([]byte{c, c})
           }
       }
       // headChar, headRunLen, firstDiff
       if t > 0 {
           headChar[i] = c
           headRunLen[i] = t
           // firstDiff: next char after head run
           if r >= n || isEmpty[r] {
               firstDiff[i] = 0
           } else if headChar[r] != c {
               firstDiff[i] = headChar[r]
           } else {
               firstDiff[i] = firstDiff[r]
           }
       } else {
           // t == 0
           headChar[i] = headChar[r]
           headRunLen[i] = headRunLen[r]
           firstDiff[i] = firstDiff[r]
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < n; i++ {
       L := length[i]
       // print length
       fmt.Fprint(writer, L, " ")
       if L <= 10 {
           // full
           fmt.Fprint(writer, pref[i])
       } else {
           // pref 5 + ... + suff 2
           fmt.Fprint(writer, pref[i], "...", suff[i])
       }
       fmt.Fprint(writer, '\n')
   }
}
