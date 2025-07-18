#include <bits/stdc++.h>
using namespace std;

int t;

int main() {
  ios::sync_with_stdio(false); cin.tie(NULL); cout.tie(NULL);
  cin >> t;
  cout.precision(20);
  while (t--) {
    long long px, py, ax, ay, bx, by;
    long long ans = INT_MAX;
    cin >> px >> py >> ax >> ay >> bx >> by;
    ans = min(
      ans,
      max(
        (ax - px) * (ax - px) + (ay - py) * (ay - py),
        ax * ax + ay * ay
      )
    );
    ans = min(
      ans,
      max(
        (bx - px) * (bx - px) + (by - py) * (by - py),
        bx * bx + by * by
      )
    );
    long double actualAns = sqrt(static_cast<long double>(ans));
    actualAns = min(
      actualAns,
      max(
        sqrt(static_cast<long double>((ax - bx) * (ax - bx) + (ay - by) * (ay - by))) / 2,
        max(
          sqrt(static_cast<long double>(ax * ax + ay * ay)),
          sqrt(static_cast<long double>((bx - px) * (bx - px) + (by - py) * (by - py)))
        )
      )
    );
    actualAns = min(
      actualAns,
      max(
        sqrt(static_cast<long double>((ax - bx) * (ax - bx) + (ay - by) * (ay - by))) / 2,
        max(
          sqrt(static_cast<long double>(bx * bx + by * by)),
          sqrt(static_cast<long double>((ax - px) * (ax - px) + (ay - py) * (ay - py)))
        )
      )
    );
    cout << actualAns << "\n";
  }
  return 0;
}