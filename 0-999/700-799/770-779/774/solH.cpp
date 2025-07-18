#include <iostream>
#include <fstream>

#include <stack>
#include <queue>
#include <deque>
#include <vector>

#include <set>
#include <unordered_set>
#include <map>
#include <unordered_map>

#include <cmath>

#include <cassert>

#include <algorithm>

using namespace std;

int main()
{
  cin.tie(nullptr);
  ios_base::sync_with_stdio(false);

  cout.setf(ios_base::fixed);
  cout.precision(28);


  int n;
  cin >> n;
  vector<int> v(n);
  for (auto& x : v) {
    cin >> x;
  }
  string S;
  int c = 0;
  int a = 0;
  for (int i = n - 1; i > 0; --i) {
    if (v[i] == 0) {
      continue;
    }
    for (int j = 0; j <= i; ++j) {
      S += (a % 2 ? 'a' : 'b');
    }
    ++a;
    for (int j = i - 1; j >= 0; --j) {
      v[j] -= i - j + 1;
    }
    --v[i];
    if (v[i] > 0) {
      ++i;
    }
  }
  while (S.length() < n) {
    S += (a % 2 ? 'a' : 'b');
    ++a;
  }
  cout << S;
  return 0;
}

/*
7
7 3 1 0 0 0 0
7
7 4 2 0 0 0 0
 */