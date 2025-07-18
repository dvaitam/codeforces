#include <iostream>
#include <vector>
#include <map>
#include <set>
#include <stack>
#include <queue>
#include <limits>
#include <string>
#include <array>
#include <sstream>
#include <numeric>
#include <list>
#include <deque>
#include <cstring>
#include <cmath>
#include <cstdlib>
#include <cstdio>
#include <random>
#include <algorithm>
#include <climits>
#include <complex>
#include <cassert>
#include <bitset>

using namespace std;

#define itr(...) begin(__VA_ARGS__), end(__VA_ARGS__)
#define fastIO() ios::sync_with_stdio(false),cin.tie(nullptr),cout.tie(nullptr)
#define debug(x...) do { cout << "\033[32;1m" << #x <<" -> "; rd_debug(x); } while (0)
void rd_debug() { cout << "\033[39;0m" << endl; }
template<class T, class... Ts> void rd_debug(const T& arg,const Ts&... args) { cout << arg << " "; rd_debug(args...); }
#define pb push_back
#define PF(x) ((x) * (x))
#define LF(x) ((x) * PF(x))

typedef long long ll;
typedef unsigned long long ull;
typedef pair<int, int> PII;
typedef pair<ll, ll> PLL;

const double eps = 1e-6;
const int MOD = 1e9 + 7;
const int inf = 0x3f3f3f3f;
const ll infl = 0x3f3f3f3f3f3f3f3fll;

int main() {
  fastIO();
  int t;
  cin >> t;
  while(t--){
    string s;
    cin >> s;
    int n;
    cin >> n;
    vector<PII> ss(s.length(), {inf, inf});
    vector<string> arr(n);
    for(int i = 0; i < n; ++i){
      cin >> arr[i];
      for(int j = 0; j < s.length(); ++j){
        if(s.substr(j, arr[i].length()) == arr[i]){
          int rp = j + (int)arr[i].length() - 1;
          int idx = i;
          if(ss[rp].second > j){
            ss[rp] = {idx, j};
          }
        }
      }
    }
    vector<PII> ans;
    int rpos = (int)s.length() - 1;
    while(rpos >= 0){
      int mv = rpos;
      int cs = -1;
      for(int i = rpos; i < s.length(); ++i){
        if(ss[i].second - 1 < mv){
          mv = ss[i].second - 1;
          cs = ss[i].first;
        }
      }
      if(mv >= rpos) break;
      else{
        ans.push_back({cs, mv + 1});
      }
      rpos = mv;
    }
    if(rpos >= 0) cout << -1 << "\n";
    else {
      cout << ans.size() << "\n";
      for(auto& it: ans){
        cout << it.first + 1 << " " << it.second + 1 << "\n";
      }
    }
  }
  return 0;
}