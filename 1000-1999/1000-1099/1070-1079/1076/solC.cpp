#include <bits/stdc++.h>
#define ll long long
#define REP(i, n) for (ll i = 0, max_i = (n); i < max_i; i++)
#define REPI(i, a, b) for (ll i = (a), max_i = (b); i < max_i; i++)
#define ALL(obj) (obj).begin(), (obj).end()
#define RALL(obj) (obj).rbegin(), (obj).rend()
#define fi first
#define se second
#define pb push_back
#define debug(x) cerr << #x << ": " << (x) << endl
#define debug2(x, y) cerr << #x << ": " << (x) << " " << #y << ": " << y << endl;
#define int long long
using namespace std;
using II = pair<int, int>;
using VII = vector<II>;
using VI = vector<int>;
using VVI = vector<VI>;
using VVVI = vector<VVI>;
template <class T = int> inline T in() { T x; cin >> x; return x; }
template <class T = int> inline bool chmax(T &a, const T &b) { if (a < b) { a = b; return true; } return false; }
template <class T = int> inline bool chmin(T &a, const T &b) { if (a > b) { a = b; return true; } return false; }
template <class T> ostream& operator<<(ostream &s, const vector<T>& d) { int n = d.size(); REP (i, n) s << d[i] << " "; return s; }
template <class T> ostream& operator<<(ostream &s, const vector<vector<T>>& dd) { for (vector<T> d: dd) s << d << endl; return s; }
struct Fast { Fast() { cin.tie(0); ios::sync_with_stdio(false); } } fast;
const int MOD = 1e9 + 7;

void solve() {
  int d = in();
  int D = d * d - 4 * d;
  if (D < 0) {
    printf("N\n"); return;
  }
  double a = (d + sqrt(D)) / 2.0;
  double b = (d - sqrt(D)) / 2.0;
  printf("Y %.15f %.15f\n", a, b);
}
signed main() {
  int t = in();
  while (t--) {
    solve();
  }
}