#include <bits/stdc++.h>

#define each(i, c) for (auto& i : c)
#define unless(cond) if (!(cond))

using namespace std;

typedef long long int lli;
typedef unsigned long long ull;
typedef complex<double> point;

template<typename P, typename Q> ostream& operator << (ostream& os, pair<P, Q> p) { os << "(" << p.first << "," << p.second << ")"; return os; }
template<typename P, typename Q> istream& operator >> (istream& is, pair<P, Q>& p) { is >> p.first >> p.second; return is; }
template<typename T> ostream& operator << (ostream& os, vector<T> v) { os << "("; each (i, v) os << i << ","; os << ")"; return os; }
template<typename T> istream& operator >> (istream& is, vector<T>& v) { each (i, v) is >> i; return is; }

template<typename T> inline T setmax(T& a, T b) { return a = std::max(a, b); }
template<typename T> inline T setmin(T& a, T b) { return a = std::min(a, b); }

int main(int argc, char *argv[])
{
  ios_base::sync_with_stdio(0);
  cin.tie(0);

  int n;
  string s;
  while (cin >> n >> s) {
    for (int i = 0; i + 1 < s.size(); ++i) {
      if (s[i] > s[i + 1]) {
        s[i] = '@';
        break;
      }
    }
    if (count(s.begin(), s.end(), '@') == 0) {
      for (int i = 0; i + 1 < s.size(); ++i) {
        cout << s[i];
      }
      cout << endl;
    } else {
      each (c, s) {
        if (c != '@') cout << c;
      }
      cout << endl;
    }
    
  }

  
  return 0;
}