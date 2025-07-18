#include <iostream>
#include <vector>
using namespace std;

int f(int i, int avoid) {
  int ans = i & -i;
  if (ans >= avoid) ans <<= 1;
  return ans;
}

int main() {
  int n, x; cin >> n >> x;
  if (x >= (1 << n)) {
    cout << (1 << n) - 1 << endl << 1;
    for (int i = 2; i < (1 << n); i++) cout << " " << (i & -i);
    return 0;
  }
  int l = 1 << (n - 1); l--;
  cout << l << endl; if (l == 0) return 0;
  int avoid = x & -x;
  cout << f(1, avoid);
  for (int i = 2; i <= l; i++) cout << " " << f(i, avoid);
  cout << endl;
  return 0;
}