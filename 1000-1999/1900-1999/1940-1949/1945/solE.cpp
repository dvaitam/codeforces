/*
  ! Though life is hard, I want it to be boiling.
  ! Created: 2024/03/19 15:53:39
*/
#include <bits/stdc++.h>
using namespace std;

#define x first
#define y second
#define int long long
#define mp(x, y) make_pair(x, y)
#define eb(...) emplace_back(__VA_ARGS__)
#define fro(i, x, y) for(int i = (x); i <= (y); i++)
#define pre(i, x, y) for(int i = (x); i >= (y); i--)
inline void JYFILE19();

typedef long long i64;
typedef pair<int, int> PII;

bool ST;
const int N = 1e6 + 10;
const int mod = 998244353;

int n, m, p[N], vs[N];

inline void solve() {
  cin >> n >> m;
  fro(i, 1, n) cin >> p[i], vs[i] = 0;
  int l = 1, r = n + 1;
  while (r - l > 1) {
    int mid = (l + r) >> 1;
    vs[mid] = 1;
    if (p[mid] <= m) {
      l = mid;
    }
    else {
      r = mid;
    }
  }
  if (p[l] == m) {
    cout << 0 <<"\n";
  }
  else {
    int id = 0;
    fro(i, 1, n) if (p[i] == m) id = i;
    if (p[l] <= m) {
      cout << 1 << "\n";
      cout << id << " " << l << "\n";
    }
    if (p[l] > m) {
      cout << 2 << "\n";
      cout << id << " " << l << "\n";
      swap(p[id], p[l]);
      id = l, l = 1, r = n + 1;
      while (r - l > 1) {
        int mid = (l + r) >> 1;
        if (p[mid] <= m) {
          l = mid;
        }
        else {
          r = mid;
        }
      }
      cout << id << " " << l << "\n";
    }
  }
}

signed main() {
  JYFILE19();
  int t; cin >> t;
  while (t--) {
    solve();
  }
  return 0;
}

bool ED;
inline void JYFILE19() {
  // freopen("", "r", stdin);
  // freopen("", "w", stdout);
  srand(random_device{}());
  ios::sync_with_stdio(0), cin.tie(0);
  double MIB = fabs((&ED-&ST)/1048576.), LIM = 32;
  cerr << "MEMORY: " << MIB << endl, assert(MIB<=LIM);
}