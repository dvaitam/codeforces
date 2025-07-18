#include <string.h>

#include <algorithm>

#include <bitset>

#include <cassert>

#include <cmath>

#include <complex>

#include <cstdio>

#include <cstdlib>

#include <cstring>

#include <ctime>

#include <deque>

#include <fstream>

#include <functional>

#include <iomanip>

#include <iostream>

#include <map>

#include <queue>

#include <set>

#include <sstream>

#include <stack>

#include <string>

// #include <unordered_map>

// #include <unordered_set>

#include <utility>

#include <vector>



#define pb push_back

#define pf push_front

#define mp make_pair



using namespace std;

typedef long long LL;

typedef unsigned long long ULL;

typedef pair<int, int> pii;

typedef pair<LL, LL> pll;



inline void EnableFileIO(const string &fileName) {

  if (fileName.empty()) return;

#ifdef ONLINE_JUDGE

  freopen((fileName + ".in").c_str(), "r", stdin);

  freopen((fileName + ".out").c_str(), "w", stdout);

#endif

  return;

}



const int INF = (1 << 30) - 1;

const LL LINF = (1LL << 61) - 1;

const double EPS = 1e-10;

const int N = 200;



int r1, c1, r2, c2;

int n, size[N];

char field[N][N];



inline void next(int &x, int &y, int &dx) {

  int nx = x + dx;



  if (nx < 0) {

    dx = 1, x = 0, y++;

  } else if (y < r1 && nx >= c1) {

    dx = -1, x = c1 - 1, y++;

  } else if (y >= r1 && nx >= c2) {

    dx = -1, x = c2 - 1, y++;

  } else {

    x += dx;

  }

}



int main() {

  ios::sync_with_stdio(false);

  cin.tie(0);

  // printf("Hello, world!\n");

  EnableFileIO("");



  cin >> r1 >> c1 >> r2 >> c2;

  cin >> n;

  for (int i = 0; i < n; i++) {

    cin >> size[i];

  }



  int x = 0, y = 0;

  int dx = 1;



  if (r1 & 1) {

    x = c1 - 1, dx = -1;

  }



  for (int i = 0; i < N; i++) {

    for (int j = 0; j < N; j++) field[i][j] = '.';

  }



  for (int i = 0; i < n; i++) {

    while (size[i]) {

      field[x][y] = i + 'a';

      size[i]--;

      next(x, y, dx);

    }

  }



  cout << "YES" << endl;

  for (int i = 0; i < max(c1, c2); i++) {

    for (int j = 0; j < r1 + r2; j++) {

      cout << field[i][j];

    }

    cout << endl;

  }



  return 0;

}