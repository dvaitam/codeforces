#include <bits/stdc++.h>
using namespace std;
// clang-format off
template <typename A, typename B> ostream& operator<<(ostream& os, const pair<A, B>& p) { return os << '(' << p.first << ' ' << p.second << ')'; }
template <typename X, typename T = typename enable_if<!is_same<X, string>::value, typename X::value_type>::type> ostream& operator<<(ostream& o, const X& v) { string s; for (const T& x : v) o << s << x, s = " "; return o; }
void deb() { cout << "\n"; }
template <typename Head, typename... Tail> void deb(Head H, Tail... T) { cout << H; if (sizeof...(T) > 0) cout << ' '; deb(T...); }
#ifdef LOCAL
#define dbg(...) cout << "[" << #__VA_ARGS__ << "]:", deb(__VA_ARGS__)
#else
#define dbg(...)
#endif
// clang-format on

void solve() {
  int n, q, l, r;
  cin >> n;
  vector<int> a(n);
  for (int& x : a) cin >> x;
  vector<pair<int, int>> b;
  for (int i = 0; i < n; i++) {
    int j = i;
    while (j < n and a[i] == a[j]) j++;
    j--;
    b.push_back({ i + 1, j + 1 });
    i = j;
  }
  cin >> q;
  while (q--) {
    cin >> l >> r;
    auto it = upper_bound(b.begin(), b.end(), make_pair(l, n + 10));
    it--;
    if (it->second < r) {
      cout << l << ' ' << it->second + 1 << '\n';
    } else {
      cout << -1 << ' ' << -1 << '\n';
    }
  }
}

int32_t main() {
  ios_base::sync_with_stdio(false);
  cin.tie(NULL);
  int T = 1;
  cin >> T;
  while (T--) solve();
  return 0;
}