#include <bits/stdc++.h>
using namespace std;
#define rep(i, a, n) for (int i = a; i < n; i++)
#define per(i, a, n) for (int i = n - 1; i >= a; i--)
#define pb push_back
#define mp make_pair
#define clr(x) memset(x, 0, sizeof(x))
#define all(x) (x).begin(), (x).end()
#define fi first
#define se second
#define sz(x) ((int)(x).size())
typedef vector<int> vi;
typedef long long ll;
typedef pair<int, int> pii;
const double PI = acos(-1.0);
const ll MOD = 1000000007;
const int LOG = 61;
struct LinearRec {
  int n;
  vector<int> first, trans;
  vector<vector<int>> bin;

  vector<int> add(vector<int> &a, vector<int> &b) {
    vector<int> result(n * 2 + 1, 0);
    for (int i = 0; i <= n; ++i) {
      for (int j = 0; j <= n; ++j) {
        if ((result[i + j] += (long long)a[i] * b[j] % MOD) >= MOD) {
          result[i + j] -= MOD;
        }
      }
    }
    for (int i = 2 * n; i > n; --i) {
      for (int j = 0; j < n; ++j) {
        if ((result[i - 1 - j] += (long long)result[i] * trans[j] % MOD) >=
            MOD) {
          result[i - 1 - j] -= MOD;
        }
      }
      result[i] = 0;
    }
    result.erase(result.begin() + n + 1, result.end());
    return result;
  }

  LinearRec(vector<int> &first, vector<int> &trans)
      : first(first), trans(trans) {
    n = first.size();
    vector<int> a(n + 1, 0);
    a[1] = 1;
    bin.push_back(a);
    for (int i = 1; i < LOG; ++i) {
      bin.push_back(add(bin[i - 1], bin[i - 1]));
    }
  }

  int calc(ll k) {
    vector<int> a(n + 1, 0);
    a[0] = 1;
    for (int i = 0; i < LOG; ++i) {
      if (k >> i & 1) {
        a = add(a, bin[i]);
      }
    }
    int ret = 0;
    for (int i = 0; i < n; ++i) {
      if ((ret += (long long)a[i + 1] * first[i] % MOD) >= MOD) {
        ret -= MOD;
      }
    }
    return ret;
  }
};
int main() {
  int n;
  ll k;
  cin >> k >> n;
  k++;
  vector<int> a(n, 0);
  a[0] = 1;
  a[n - 1] = 1;
  vector<int> b(n, 1);
  LinearRec f(b, a);
  printf("%d\n", f.calc(k));
  return 0;
}