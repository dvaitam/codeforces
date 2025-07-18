#define _USE_MATH_DEFINES
#include <cstdio>
#include <cstring>
#include <iostream>
#include <string>
#include <sstream>
#include <vector>
#include <queue>
#include <list>
#include <set>
#include <map>
#include <unordered_set>
#include <unordered_map>
#include <algorithm>
#include <complex>
#include <cmath>
#include <numeric>
#include <bitset>

using namespace std;

typedef long long int64;
typedef pair<int, int> ii;
const int INF = 1 << 30;

int main() {
  double ret = 0.0;
  for (int i = 3; i <= 5; ++i) {
    double len;
    cin >> len;
    double r = len / 2 / sin(M_PI / i);
    double h = sqrt(len * len - r * r);
    double area = r * r * sin(2 * M_PI / i) / 2 * i;
    ret += area * h / 3;
  }
  printf("%.10f\n", ret);
  return 0;
}