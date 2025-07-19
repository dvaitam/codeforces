package main

import (
   "bufio"
   "fmt"
   "os"
)

const AwA = 1000010

var trans [AwA][26]int32
var fail [AwA]int32
var lengthArr [AwA]int32
var pos [AwA]int32
var lstpos [AwA]int32
var g [AwA]int32
var fArr [AwA]int32

var tot int32

func initSAM() {
   tot = 1
   pos[0] = 1
   fail[1] = 0
   lengthArr[1] = 0
}

func insertSAM(cnt, c int32) {
   cur := tot + 1
   tot = cur
   lst := pos[cnt-1]
   pos[cnt] = cur
   lengthArr[cur] = cnt

   for lst != 0 && trans[lst][c] == 0 {
       trans[lst][c] = cur
       lst = fail[lst]
   }
   if lst == 0 {
       fail[cur] = 1
       return
   }
   p := trans[lst][c]
   if lengthArr[p] != lengthArr[lst]+1 {
       q := tot + 1
       tot = q
       lengthArr[q] = lengthArr[lst] + 1
       fail[q] = fail[p]
       // copy transitions
       trans[q] = trans[p]
       for lst != 0 && trans[lst][c] == p {
           trans[lst][c] = q
           lst = fail[lst]
       }
       fail[p] = q
       fail[cur] = q
   } else {
       fail[cur] = p
   }
}

func moveSAM(u *int32, length int32) {
   for *u != 1 && lengthArr[fail[*u]] >= length {
       *u = fail[*u]
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int32
   fmt.Fscan(reader, &n)
   var str string
   fmt.Fscan(reader, &str)
   s := []byte(str)
   // reverse string
   for i, j := 0, int(n)-1; i < j; i, j = i+1, j-1 {
       s[i], s[j] = s[j], s[i]
   }

   initSAM()
   for i := int32(1); i <= n; i++ {
       insertSAM(i, int32(s[i-1]-'a'))
   }

   var lst int32 = 1
   for i := int32(1); i <= n; i++ {
       c := int32(s[i-1] - 'a')
       cur := trans[lst][c]
       curfa := cur
       fArr[i] = fArr[i-1] + 1
       moveSAM(&curfa, fArr[i]-1)
       for fArr[i] != 1 {
           if g[lst] >= fArr[i]-1 || g[curfa] >= fArr[i]-1 {
               break
           }
           fArr[i]--
           moveSAM(&cur, fArr[i])
           moveSAM(&lst, fArr[i]-1)
           moveSAM(&curfa, fArr[i]-1)
           p := i - fArr[i]
           u := lstpos[p]
           if g[u] < fArr[p] {
               g[u] = fArr[p]
           }
           u = fail[u]
           for u != 0 && g[u] < lengthArr[u] {
               g[u] = lengthArr[u]
               u = fail[u]
           }
       }
       lst = cur
       lstpos[i] = cur
   }

   var ans int32
   for i := int32(1); i <= n; i++ {
       if fArr[i] > ans {
           ans = fArr[i]
       }
   }
   fmt.Println(ans)
}
