package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const (
   seed uint64 = 300 + 57
   base uint64 = 10000 + 317
)

func hsv(s string) uint64 {
   var h uint64
   for i := 0; i < len(s); i++ {
       h = h*seed + uint64(s[i])
   }
   return h
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // map host to set of paths
   hostPaths := make(map[string]map[string]struct{}, n)
   for i := 0; i < n; i++ {
       var url string
       fmt.Fscan(in, &url)
       // skip "http://"
       s := url
       if len(s) > 7 {
           s = s[7:]
       }
       // parse host
       var j int
       for j < len(s) {
           c := s[j]
           if (c >= 'a' && c <= 'z') || c == '.' {
               j++
           } else {
               break
           }
       }
       host := s[:j]
       // parse path
       k := j
       for k < len(s) {
           c := s[k]
           if (c >= 'a' && c <= 'z') || c == '.' || c == '/' {
               k++
           } else {
               break
           }
       }
       path := s[j:k]
       if len(path) == 0 {
           path = "$"
       }
       m, ok := hostPaths[host]
       if !ok {
           m = make(map[string]struct{})
           hostPaths[host] = m
       }
       m[path] = struct{}{}
   }
   // compute hash for each host
   groups := make(map[uint64][]string)
   for host, pathsMap := range hostPaths {
       // collect and sort paths
       paths := make([]string, 0, len(pathsMap))
       for p := range pathsMap {
           paths = append(paths, p)
       }
       sort.Strings(paths)
       // compute combined hash
       var hv uint64
       for _, p := range paths {
           hv = hv*base + hsv(p)
       }
       groups[hv] = append(groups[hv], host)
   }
   // collect duplicate groups
   var keys []uint64
   for hv, hs := range groups {
       if len(hs) > 1 {
           keys = append(keys, hv)
       }
   }
   sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
   // output
   fmt.Println(len(keys))
   for _, hv := range keys {
       hs := groups[hv]
       sort.Strings(hs)
       // print hosts
       for i, host := range hs {
           if i > 0 {
               fmt.Print(" ")
           }
           fmt.Print("http://", host)
       }
       fmt.Println()
   }
}
