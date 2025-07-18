#include <cstdio>
#include <cstring>
#include <algorithm>
using namespace std;

int main(void) {
  int n, m;
  scanf("%d%d", &n, &m);
  double ans = 1e18;
  while (n--) {
    int a, b;
    scanf("%d%d", &a, &b);
    ans = min(ans, (double)m * a / b);
  }
  printf("%.20lf\n", ans);
  return 0;
}