#include <algorithm>
#include <iostream>
#include <cmath>
#include <numeric>
#include <vector>

int main() {
  std::cin.tie(0)->sync_with_stdio(0);

  int n;
  std::cin >> n;

  long long x1, y1, x2, y2;
  std::cin >> x1 >> y1 >> x2 >> y2;

  long long cross = x1 * y2 - x2 * y1;

  if (cross < 0) {
    std::swap(x1, x2);
    std::swap(y1, y2);
    cross = -cross;
  }

  if (n != cross) return std::cout << "NO\n", 0;

  long long dx = std::abs((y1 * x2 - x1 * y2) / std::gcd(y1, y2));
  long long dy = std::abs((y1 * x2 - x1 * y2) / std::gcd(x1, x2));

  std::cout << "YES\n";
  for (int i = 0; i < dx; ++i)
    for (int j = 0; j < dy; ++j) {
      std::cout << i << ' ' << j << '\n';
      if (!--n) return 0;
    }
}