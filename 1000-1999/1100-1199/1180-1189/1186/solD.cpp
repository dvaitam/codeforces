#include <bits/stdc++.h>

#define fi first
#define se second
using namespace std;
#define inf (1e10)
const int mod = 100003;
typedef long long ll;
typedef pair<int, int> pii;
typedef pair<ll, ll> pll;
const int maxn = 100;
const double eps = 1e-9;
int cmp(double x) {
	if (fabs(x) <= eps) return 0;
	return x < 0 ? -1 : 1;
}
struct point {
		double x, y;
		point() {}
		point(int x1, int y1) { x = x1, y = y1; }
		point operator-(const point &a) const {
			return point(x - a.x, y - a.y);
		}
		point operator+(const point &a) const {
			return point(x + a.x, y + a.y);
		}
		point operator*(const double &r) const {
			return point(x * r, y * r);
		}
		double norm() {
			return sqrt(x * x + y * y);
		}
} po[maxn];

double det(point a, point b) {
	return a.x * b.y - a.y * b.x;
}
double dot(point a, point b) {
	return a.x * a.y + b.x * b.y;
}
bool pointonseg(point p, point s, point t) {
	return cmp(det(p - s, t - s)) == 0 && cmp(dot(p - s, p - t)) <= 0;
}
double dis_point_seg(point p, point s, point t, point &ret) {
	if (cmp(dot(p - s, t - s)) < 0) {
		ret = s;
		return (p - s).norm();
	}
	if (cmp(dot(p - t, s - t)) < 0) {
		ret = t;
		return (p - t).norm();
	}
	double r = dot((t - s), p - s) / dot(t - s, t - s);
	ret = s + (t - s) * r;
	return fabs(det(s - p, t - p) / (s - t).norm());
}
struct polyon {
		point a[maxn];
		int n;
		int PointIn(point t) {
			int num = 0, i, d1, d2, k;
			a[n] = a[0];
			for (int i = 0; i < n; ++i) {
				if (pointonseg(t, a[i], a[i + 1])) return 2;
				k = cmp(det(a[i + 1] - a[i], t - a[i]));
				d1 = cmp(a[i].y - t.y);
				d2 = cmp(a[i + 1].y - t.y);
				if (k > 0 && d1 <= 0 && d2 > 0) num++;
				if (k < 0 && d2 <= 0 && d1 > 0) num--;
			}
			return num != 0;
		}
		point pointdis(point t) {
			a[n] = a[0];
			if (PointIn(t)) return t;
			point res = a[0];
			double dis = inf;
			for (int i = 0; i < n; ++i) {
				point temp;
				double d = dis_point_seg(t, a[i], a[i + 1], temp);
				if (cmp(d - dis) < 0) dis = d, res = temp;
			}
			return res;
		}
} P[maxn];
int n;
point near(point t) {
	double dis = inf;
	point ans;
	for (int i = 0; i < n; ++i) {
		point temp = P[i].pointdis(t);
		if (cmp((t - temp).norm() - dis) < 0) {
			dis = (t - temp).norm();
			ans = temp;
		}
	}
	return ans;
}
double cal(point s, point t) {
	point p1 = near(s);
	point p2 = near(t);
	double len = (p1 - p2).norm();
	double l = 0, r = len;
	for (int i = 0; i < 100; ++i) {
		double mid1, mid2;
		mid1 = l + (r - l) / 3, mid2 = l + 2 * (r - l) / 3;
		point p11, p22;
		p11 = s + (t - s) * (mid1 / len);
		p22 = s + (t - s) * (mid2 / len);
		if (max((p11 - p1).norm(), (p11 - p2).norm()) <
		    max((p22 - p1).norm(), (p22 - p2).norm()))
			l = mid1;
		else r = mid2;
	}
	point p11 = s + (t - s) * (l / len);
	return max((p11 - p1).norm(), (p11 - p2).norm());
}
ll slove(ll n, ll m) {
	ll L = m / 3;
	if (m % 3 == 2) L++;
	ll R = n / 2;
	if (n % 2 == 1) R++;
	ll ret = L * R;
	if (m - (L * 2 + L - 1) >= 2) {
		ret += n / 3;
		if (n % 3 == 2) ret++;
	}
	return ret;
}
inline char nc() {
	static char buf[100000], *p1 = buf, *p2 = buf;
	if (p1 == p2) {
		p2 = (p1 = buf) + fread(buf, 1, 100000, stdin);
		if (p1 == p2)return EOF;
	}
	return *p1++;
}
inline int read(int &x) {
	char c = nc(), b = 1;
	for (; !(c >= '0' && c <= '9'); c = nc())if (c == '-')b = -1;
	for (x = 0; c >= '0' && c <= '9'; x = x * 10 + c - '0', c = nc());
	x *= b;
	return b;
}
pii PA[(int) 1e5 + 100];
int main() {
//	freopen("1.in", "r", stdin);
//	freopen("1.out", "w", stdout);
	int n;
	read(n);
	ll sum = 0;
	for (int i = 1; i <= n; ++i) {
		int a, b;
		char s[1000];
		int op = read(a);
		read(b);
		if (a < 0 || op == -1) b *= -1;
		sum += b;
		PA[i] = {a, b};
	}
	ll res = sum / (int) 1e5;
	assert(sum % (int) 1e5 == 0);
	for (int i = 1; i <= n; ++i) {
		if (res != 0 && PA[i].se != 0) {
			if (res > 0) {
				if (PA[i].se > 0) {
					printf("%d\n", PA[i].fi + 1);
					res--;
				} else {
					printf("%d\n", PA[i].fi);
				}
			} else {
				if (PA[i].se < 0) {
					printf("%d\n", PA[i].fi - 1);
					res++;
				} else {
					printf("%d\n", PA[i].fi);
				}
			}
		} else printf("%d\n", PA[i].fi);
	}
	assert(res == 0);
	return 0;
}