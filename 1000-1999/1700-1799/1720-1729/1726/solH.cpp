#include <bits/stdc++.h>

#define sz(x) (int)(x).size()

#define all(x) (x).begin(), (x).end()

#define st first

#define nd second

using namespace std;

typedef long long ll;

typedef long double ld;

typedef pair<ll, ll> pt;

const int mod = 1e9 + 7;

const int N = 1e6 + 5;

const ld pi = acos(-1);

const ld eps = 1e-8;

pt operator-(pt p, pt q) {

  return {p.st - q.st, p.nd - q.nd};

}

ll operator*(pt p, pt q) {

  return p.st * q.st + p.nd * q.nd;

}

ld len(pt p) {

  return sqrt(p * p);

}

ld angle(pt p, pt q, pt r) {

  return acos((p - q) * (r - q) / (len(p - q) * len(r - q)));

}

ld area(ld a) {

  if (abs(a - pi / 2) < eps)

    return 3 * pi / 32. - (a - pi / 2.) / 4.;

  if (abs(a - pi) < 1e-4)

    return (pi - a) / 6;

  return (tan(pi / 2 - a) + (pi - a) * (2 + 1. / (sin(a) * sin(a))) - 4 * sin(2 * a)) / 16. + sin(2 * a) / 4;

}

ld xx(ld a, ld t) {

  return (cos(t) * sin(a + t) * sin(a + t) - cos(a) * sin(t) * sin(t) * cos(a + t)) / (sin(a) * sin(a));

}

ld yy(ld a, ld t) {

  return -sin(t) * sin(t) * cos(a + t) / sin(a);

}

// solve xx(a, t) = x

ld calc(ld a, ld x) {

  ld lo = 0, hi = pi - a;

  while (hi - lo > eps) {

    ld mi = (lo + hi) / 2.;

    if (xx(a, mi) > x)

      lo = mi;

    else

      hi = mi;

  }

  return lo;

}

ld integrate(ld a, ld t) {

  ld r = -8 * t + 2 * sin(4 * t) + cos(2 * a) * (4 * t - 3 * sin(2 * t) - sin(4 * t) + sin(6 * t)) + 8 * (3 + 2 * cos(2 * t)) * pow(sin(t), 4) * sin(2 * a);

  return -r / (64 * sin(a) * sin(a));

}

int main() {

  ios_base::sync_with_stdio(false); cin.tie(0);

  int n; cin >> n;

  vector<pt> a(n);

  for (int i = 0; i < n; i++)

    cin >> a[i].st >> a[i].nd;

  a.push_back(a[0]);

  a.push_back(a[1]);

  a.push_back(a[2]);



  if (n == 4) {

    for (int i = 0; i < 2; i++) {

      if (a[i].st == a[i + 1].st && a[i + 2].st == a[i + 3].st && (abs(a[i + 2].st - a[i + 1].st) == 1 || abs(a[i].nd - a[i + 1].nd) == 1)) {

        cout << abs(a[i + 2].st - a[i + 1].st) * abs(a[i].nd - a[i + 1].nd) << endl;

        return 0;

      }

    }

  }



  ld ans = 0;

  for (int i = 1; i <= n; i++) {

    ld al = angle(a[i - 1], a[i], a[i + 1]);

    ans += area(al);

  }

  for (int i = 1; i <= n; i++) {

    ld d = len(a[i] - a[i + 1]);

    if (d >= 2 - eps)

      continue;

    ld al = angle(a[i - 1], a[i], a[i + 1]);

    ld be = angle(a[i], a[i + 1], a[i + 2]);

    ld lo = 0, hi = pi - al;

    while (hi - lo > eps) {

      ld mi = (lo + hi) / 2.;

      if (d - xx(al, mi) <= 1 && yy(al, mi) < yy(be, calc(be, d - xx(al, mi))))

        lo = mi;

      else

        hi = mi;

    }

    ans -= integrate(al, lo);

    ans -= integrate(be, calc(be, d - xx(al, lo)));

  }

  cout << fixed << setprecision(15) << ans << endl;

}