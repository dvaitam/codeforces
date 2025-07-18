#include <bits/stdc++.h>
#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>

using namespace std;
using namespace __gnu_pbds;

#define PB push_back
#define EB emplace_back
#define MP make_pair
#define F first
#define S second
#define FBO find_by_order
#define OOK order_of_key
#define MSET(m,v) memset((m),(v),sizeof(m))
#define UB(c,v) (upper_bound((c).begin(),(c).end(),(v)))
#define LB(c,v) (lower_bound((c).begin(),(c).end(),(v)))
#define INREP int TC;cin>>TC;while(TC--)
#define CPPIO ios::sync_with_stdio(0);cin.tie(0);cout.tie(0)
#define PI 3.1415926535897932384626433832795028841971693993751
#define TRACE
#ifdef TRACE
#define trace(...) __f(#__VA_ARGS__, __VA_ARGS__)
template <typename Arg1>
void __f(const char* name, Arg1&& arg1) {
  cerr << name << " : " << arg1 << endl;
}
template <typename Arg1, typename... Args>
void __f(const char* names, Arg1&& arg1, Args&&... args) {
  const char* comma = strchr(names + 1, ',');
  cerr.write(names, comma - names) << " : " << arg1<<" | ";
  __f(comma+1, args...);
}
#else
#define trace(...)
#endif

typedef long long ll;
typedef long double ld;
typedef vector <int> vi;
typedef vector <ll> vll;
typedef pair <int,int> pii;
typedef pair <ll,ll> pll;
typedef vector <pii> vpii;
typedef vector <pll> vpll;
typedef tree <int, null_type, less <int>, rb_tree_tag, tree_order_statistics_node_update> pbds;

// SOLUTION

int main () {
 CPPIO;
  int n;
  cin >> n;
  string s, t;
  cin >> s >> t;
  if (!is_permutation (s.begin(), s.end(), t.begin())) return cout << -1, 0;
  vi ans;
  for(int i = 0; i < n; ++i)
  for(int j = i; j < n; ++j) if (t[n - i - 1] == s[j]) {
   ans.insert (ans.end(), {n, j, 1});
   reverse (s.begin() + j, s.end());
   rotate (s.begin(), s.begin() + n - 1, s.end());
      // trace (s);
   break;
  }
  /*for (int i = n - 1; i >= 0; i--) {
    if (s[i] == t[i]) continue;
    int j = i - 1;
    for (j; j >= 0; j--) if (t[i] == s[j]) break;
    ans.EB (n - 1 - j);
    string ns = "";
    for (int j1 = n - 1; j1 > j; j1--) ns += s[j1];
    for (int j1 = 0; j1 <= j; j1++) ns += s[j1];
    s = ns;
    reverse (s.begin(), s.end());
    ans.EB (n);
    ans.EB (n - 1 - i);    
    ns = "";
    for (int j1 = n - 1; j1 > i; j1--) ns += s[j1];
    for (int j1 = 0; j1 <= i; j1++) ns += s[j1];
    s = ns;
    reverse (s.begin(), s.end());
    ans.EB (n);     
  }*/
  cout << ans.size() << "\n";
  for (auto i : ans) cout << i << " ";
 return 0;
}