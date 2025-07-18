#include "bits/stdc++.h"

using namespace std;

#define ll long long
const int ALPH = 26;

int main() {
  ios_base::sync_with_stdio(false);
  cin.tie(NULL);
  
  int t; cin >> t;
  while (t--) {

    ll n, m, k; cin >> n >> m >> k;

    string s; cin >> s;

    vector<int> cnt(ALPH);
    for (char c : s) {
      cnt[c - 'A']++;
    }

    vector<int> dp(k + 1);
    dp[0] = 1;
    for (int x : cnt) {
      if (!x) continue;
      for (int j = k; j >= x; j--) {
        dp[j] += dp[j - x];
      }
    }

    ll ans = k * k;

    for (int x : cnt) {
      if (!x) continue;
      for (int j = x; j <= k; j++) {
        dp[j] -= dp[j - x];
      }

      for (int i = 0; i <= k; i++) {
        if (!dp[i]) continue;

        int rm = k - i - x;

        ll n1 = n - min(n, (ll)i);
        ll m1 = m - min(m, (ll)rm);

        if (n1 + m1 > x) continue;
        ans = min(ans, n1 * m1);

      }

      for (int j = k; j >= x; j--) {
        dp[j] += dp[j - x];
      }
    }

    cout << ans << "\n";

  }
  
}