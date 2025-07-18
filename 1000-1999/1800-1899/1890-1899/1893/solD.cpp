#include <bits/stdc++.h>
 
#define pb push_back
#define int long long
#define ld long double
typedef long long ll;
#define all(x) x.begin(), (x).end()
#define rall(x) x.rbegin(), (x).rend()
using namespace std;
const int M = 998'244'353;
 
void solve() {
  int n, m;
  cin >> n >> m;
  vector<int> a(n);
  for (int i = 0; i < n; i++) {
    cin >> a[i];
  }
  vector<int> s(m), d(m);
  for (int i = 0; i < m; i++) {
    cin >> s[i];
  }
  for (int i = 0; i < m; i++) {
    cin >> d[i];
  }
  vector<int> cnt(n + 1);
  for (int i = 0; i < n; i++) {
    cnt[a[i]]++;
  }
 
  set<pair<int, int>> cubes;
  for (int x = 1; x <= n; x++) {
    if (cnt[x] == 0) continue;
    cubes.insert({cnt[x], x});
  }
 
  vector<vector<int>> ans(m);
 
  for (int i = 0; i < m; i++) {
    ans[i].assign(s[i], 0);
    for (int j = 0; j < s[i]; j++) {
      if (j >= d[i]) {
        if (cnt[ans[i][j - d[i]]] > 0) {
          cubes.insert({cnt[ans[i][j - d[i]]], ans[i][j - d[i]]});
        }
      }
      if (cubes.empty()) {
        cout << "-1\n";
        return;
      }
      ans[i][j] = (*cubes.rbegin()).second;
      cubes.erase(*cubes.rbegin());
      cnt[ans[i][j]]--;
    }
    for (int j = s[i]; j < s[i] + d[i]; j++) {
      if (cnt[ans[i][j - d[i]]] > 0) {
        cubes.insert({cnt[ans[i][j - d[i]]], ans[i][j - d[i]]});
      }
    }
  }
 
  for (int i = 0; i < m; i++) {
    for (int j = 0; j < s[i]; j++) {
      cout << ans[i][j] << ' ';
    }
    cout << '\n';
  }
 
}
 
 
signed main() {
  ios::sync_with_stdio(false);
  cin.tie(nullptr);
  cout.tie(nullptr);
  int t;
  cin >> t;
  while (t--) {
    solve();
  }
}