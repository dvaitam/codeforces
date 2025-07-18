#include <bits/stdc++.h>
using namespace std;

#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>
using namespace __gnu_pbds;

#define rep(i, a, b) for (int i = a; i < (b); ++i)
#define sz(x) (int)(x).size()
#define pb push_back
#define all(x) (x).begin(),(x).end()
#define rall(x) (x).rbegin(),(x).rend()
#define fi first
#define se second
#define n_l '\n'

using ll = long long;
typedef tree<ll, null_type, less<ll>, rb_tree_tag, tree_order_statistics_node_update> ordered_set;
using ld = long double;
using ull = uint64_t;
using pii = pair<int, int>;
using pll = pair<ll, ll>;
using vi = vector<int>;
using vl = vector<ll>;
using vvl = vector<vector<ll>>;

#define gg(...) [](const auto&...x){ char c='='; cerr<<#__VA_ARGS__<<" "; ((cerr<<exchange(c,',')<<" "<<x),...); cerr<<endl; }(__VA_ARGS__);
int nxt() { int x; cin >> x; return x; }


// const uint64_t mod = (1ULL << 61) - 1;
// const uint64_t seed = chrono::system_clock::now().time_since_epoch().count();
// const uint64_t base = mt19937_64(seed)() % (mod / 3) + (mod / 3);

// #include <atcoder/dsu>
//  #include <atcoder/lazysegtree>
// #include <atcoder/segtree>
// #include <atcoder/modint>
// using namespace atcoder;
// using mint = modint998244353;
const ll INF = 1e18;

ll tests;

ll n, m, p;

// mint fac[333333];

void solve() {
  cin >> n >> m >> p;

  ll av, bv;

  rep(i, 0, n) {
    ll x = nxt();
    if (x % p != 0) {
      av = i;
    }
  }
  rep(i, 0, m) {
    ll x = nxt();
    if (x % p != 0) {
      bv = i;
    }
  }

  cout << av+bv << endl;
}


int main() {
  cin.tie(nullptr);
  ios::sync_with_stdio(false);

  // cin >> tests;
  // while(tests--) {
    solve();
  // }
}