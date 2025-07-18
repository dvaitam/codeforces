#include <bits/stdc++.h>
using namespace std;
#define REP(i, n) for (int i = 0; i < (int) (n); ++i)

const int N = 1<<20;
bitset<N> isPrime;
bitset<N> diff;
int a[128];
vector<int> primes;
int n;

int main() {
  ios::sync_with_stdio(false); cin.tie(0);

  isPrime.flip();
  for (int p = 2; p < 10000; p++) if (isPrime[p]) {
    primes.push_back(p);
    for (int x = p*p; x < N; x += p) isPrime[x] = false;
  }

  int tc; cin >> tc;
  while (tc--) {
    diff.reset();
    cin >> n;
    REP(i, n) cin >> a[i];
    REP(i, n) REP(j, i) diff[abs(a[i] - a[j])] = true;

    int step = 0;
    for (int p : primes) {
      bool ok = true;
      for (int x = p; x < N; x += p) if (diff[x]) {
        ok = false;
        break;
      }
      if (ok) {
        step = p;
        break;
      }
    }

    if (step) {
      cout << "YES\n";
      REP(i, n) cout << i*step + 1 << " \n"[i+1==n];
    } else cout << "NO\n";
  }
  return 0;
}