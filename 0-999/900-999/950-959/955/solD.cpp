#include <iostream>
#include <cstdint>
#include <vector>
#include <string>
#include <algorithm>

std::string s;
std::string t;

uint32_t k;

std::vector<uint32_t> match_pos(const std::string &s, const std::string &pat) {
 std::vector<uint32_t> ans;
 ans.resize(pat.size() + 1);
 std::fill(ans.begin(), ans.end(), s.size());

 std::vector<uint32_t> kmpArr;
 kmpArr.reserve(pat.size() + 1);
 kmpArr.push_back(0);

 for (uint32_t i = 0;i < pat.size();i++) {
  uint32_t p = kmpArr.back();
  while ((p > 0) && (pat[p] != pat[i])) {
   p = kmpArr[p];
   
  }
  if ((pat[p] == pat[i]) && (p < i)) {
   kmpArr.push_back(p + 1);
  } else {
   kmpArr.push_back(0);
  }
 }
 uint32_t curCnt = 0;
 for (uint32_t i = 0;i < s.size();i++) {
  if (curCnt == pat.size()) {
   curCnt = kmpArr[curCnt];
  }
  while ((curCnt > 0) && (pat[curCnt] != s[i])) {
   curCnt = kmpArr[curCnt];
  }
  if (pat[curCnt] == s[i]) {
   curCnt++;
  } else {
   curCnt = 0;
  }
  if (i + 1 < k) {
   if (curCnt == pat.size()) {
    ans[curCnt] = std::min(ans[curCnt], k - 1);
   }
   continue;
  }
  uint32_t t = curCnt;
  while ((t > 0) && (ans[t] > i)) {
   ans[t] = i;
   t = kmpArr[t];
  }
  ans[t] = std::min(ans[t], i);
 }
 return ans;
}

void input() {
 std::ios::sync_with_stdio(false);
 std::cin.tie(nullptr);
 uint32_t sLen, tLen;
 std::cin >> sLen >> tLen >> k;
 s.reserve(sLen);
 t.reserve(tLen);
 std::cin >> s >> t;
}

template <class T>
std::ostream& operator<< (std::ostream &s, const std::vector<T> &v) {
 for (const T &i : v) {
  s << i << "; ";
 }
 return s;
}

int main() {
 input();
 std::vector<uint32_t> forw_arr = match_pos(s, t.substr(0, k));

 std::reverse(s.begin(), s.end());
 std::reverse(t.begin(), t.end());
 
 std::vector<uint32_t> rev_arr = match_pos(s, t.substr(0, k));

 //std::cerr << forw_arr << "\n" << rev_arr << "\n";

 for (uint32_t i = 0;i < forw_arr.size();i++) {
  uint32_t j = t.size() - i;
  if (j >= rev_arr.size()) {
   continue;
  }
  if (forw_arr[i] + rev_arr[j] + 1 < s.size()) {
   std::cout << "Yes\n";
   std::cout << forw_arr[i] - k + 2 << " ";
   std::cout << s.size() - rev_arr[j] << "\n";
   return 0;
  }
 }
 std::cout << "No\n";
}