#include <bits/stdc++.h>
#define forn(i, n) for (int i = 0; i < int(n); i++)
#define forsn(i, s, n) for (int i = s; i < int(n); i++)
#define sz(x) int(x.size())
#define all(x) x.begin(), x.end()
#define DBG(x) cerr << #x << " = " << x << endl

using namespace std;
using tint = long long;
using vi = vector<tint>;

inline void fastIO() {
  ios_base::sync_with_stdio(false);
  cin.tie(NULL);
}

template <typename T>
inline void chmax(T &lhs, T rhs) {
  lhs = max(lhs, rhs);
}

template <typename T>
inline void chmin(T &lhs, T rhs) {
  lhs = min(lhs, rhs);
}

template <typename T>
ostream &operator<<(ostream &os, vector<T> &v) {
  os << "[";
  forn(i, sz(v)) {
    if (i > 0) os << ", ";
    os << v[i];
  }
  os << "]";
  return os;
}

template <typename T, typename U>
ostream &operator<<(ostream &os, pair<T, U> &p) {
  os << "(" << p.first << ", " << p.second << ")";
  return os;
}

int main() {
  fastIO();
  int t;
  cin >> t;

  while (t--) {
    string a, b;
    cin >> a >> b;

    if (a == b) {
      cout << "YES\n" << a << '\n';
      continue;
    }

    if (a[0] == b[0]) {
      cout << "YES\n" << a[0] << "*\n";
      continue;
    }

    if (a.back() == b.back()) {
      cout << "YES\n"
           << "*" << a.back() << '\n';
      continue;
    }

    int n = sz(a);
    string found = "";
    forn(i, n - 1) {
      if (b.find(a.substr(i, 2)) != string::npos) {
        found = a.substr(i, 2);
        break;
      }
    }

    if (found.empty()) {
      cout << "NO\n";
    } else {
      cout << "YES\n*" << found << "*\n";
    }
  }
}