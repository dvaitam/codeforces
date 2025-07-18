#include <iostream>
#include <vector>
#include <map>
#include <set>
#include <unordered_map>
#include <unordered_set>
#include <cstdio>
#include <queue>
#include <stack>

#define F first
#define S second

using namespace std;

using ll = long long;
using ld = long double;
using pii = pair<int, int>;
using pll = pair<ll, ll>;

ll MOD = 1000000007;

ll get_sum(ll n, ll i, ll k = 1) {
  ll sum = 0;
  ll m = n;
  while (m > 0) {
    sum += m;
    m -= k;
    k *= i;
  }
  return sum;
}

int main() {
  ios_base::sync_with_stdio(false);
  ll n, s;
  cin >> n >> s;
  for (ll i = 1; i < n; ++i) {
    ll sum = get_sum(n, i);
    //cout << sum << endl;
    if (i == 1 && sum < s) {
      cout << "No\n";
      return 0;
    }
    if (sum > s) {
      continue;
    }
    cout << "Yes\n";
    if (i == 1) {
      cout << 1;
      for (ll j = 2; j < n; ++j) {
        cout << " " << j;
      }
      cout << endl;
      return 0;
    }
    ll k = i;
    ll m = n - 1;
    s -= n;
    ll start = 1;
    ll finish = 1;
    while (m > 0) {
      ll l = 1;
      ll r = k;
      while (l < r) {
        ll e = (l + r) / 2;
        sum = get_sum(m, i, e);
        if (sum > s) {
          l = e + 1;
        } else {
          r = e;
        }
      }
      k = l;
      //cout << get_sum(m, i, k) << " " << k << " " << s << endl;
      for (ll j = 0; j < k; ++j) {
        cout << start + (j / i) << " ";
      }
      start = finish + 1;
      finish += k;
      s -= m;
      m -= k;
      k *= i;
    }
    cout << endl;
    return 0;
  }
  cout << "No\n";
  return 0;
}