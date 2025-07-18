# include <bits/stdc++.h>
# define ll unsigned long long
# pragma optimize(2)

using namespace std;

mt19937 R(9);
const int N = 1 << 18 | 5;
unordered_map<ll, ll> mp;

ll s[N];
int n;

struct Node {
int l, r;
ll hs() {
  return (ll)l << 32 | r;
}
ll hsv() {
  return s[l] ^ s[r];
}
};

void sov() {
  cin >> n, mp.clear();
  n = 1 << n + 1;
  for (int i = 1; i <= n; i++) {
    cin >> s[i];
    s[i] ^= s[i - 1];
  }
  while (1) {
    Node q = {int(R() % (n + 1)), int(R() % (n + 1))};
    while (q.l == q.r) {
      q.l = int(R() % (n + 1));
    }
    if (q.l > q.r) {
      swap(q.l, q.r);
    }
    ll vi = q.hs(), vv = q.hsv();
    auto it = mp.find(vv);
    if (it == mp.end()) {
      mp.insert({vv, vi});
    } else {
      ll ci = it->second;
      vector<int> ans(4);
      ans[0] = vi >> 32, ans[1] = vi & (0xffffffff);
      ans[2] = ci >> 32, ans[3] = ci & (0xffffffff);
      if (ans[0] != ans[2] && ans[1] != ans[3]) {
        sort(ans.begin(), ans.end());
        cout << ans[0] + 1 << ' ' << ans[1] << ' ' << ans[2] + 1 << ' ' << ans[3] << '\n';
        return;
      } else {
        it->second = vi;
      }
    }
  }
}

int main() {
# ifndef ONLINE_JUDGE
  freopen("in.txt", "r", stdin);
# endif
  ios::sync_with_stdio(0);
  cin.tie(0), cout.tie(0);
  int T; cin >> T;
  while (T--) {
    sov();
  }
}