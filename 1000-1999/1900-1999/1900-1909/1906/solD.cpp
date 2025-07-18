#include <bits/stdc++.h>
using namespace std;
typedef long long LL;
typedef long double LD;
typedef pair <int, int> pii;
 
 
int sgn (LL x) { return x > 0 ? 1 : (x < 0 ? -1 : 0); }
LD sqr (LD x) { return x * x; }
 
template <class T = LL>
struct Point {
	T x, y;
	Point () {}
	Point (T a, T b) {x = a, y = b;}
	Point operator + (Point a) const { return {x + a.x, y + a.y}; }
	Point operator - (Point a) const { return {x - a.x, y - a.y}; }
	Point operator * (T a) const { return {x * a, y * a}; }
	Point operator / (T a) const { return {x / a, y / a}; }
	bool operator < (Point a) const { return x < a.x || (x == a.x && (y < a.y)); }
	bool operator > (Point a) const { return x > a.x || (x == a.x && (y > a.y)); }
	bool operator == (Point a) const { return x == a.x && y == a.y; }
	void read() {
		cin >> x >> y;
	}
	friend ostream & operator << (ostream &os, Point &a) {
		return os << a.x << "," << a.y;
	}
};
 
using point = Point<LL>;
using pointLD = Point<LD>;
 
#define cp const point &
 
LL det (cp a, cp b) { return a.x * b.y - a.y * b.x; }
LL dot (cp a, cp b) { return a.x * b.x + a.y * b.y; }
template <class T1, class T2>
LD dis (T1 a, T2 b) { return sqrt (sqr (a.x - b.x) + sqr (a.y - b.y)); }
 
#define cl line
struct line {
	point s, t;
};
pointLD line_inter(cl a, cl b) {
	LD s1 = det(a.t - a.s, b.s - a.s);
	LD s2 = det(a.t - a.s, b.t - a.s);
	return (pointLD(b.s.x, b.s.y) * s2 - pointLD(b.t.x, b.t.y) * s1) / (s2 - s1);
}
bool has_inter(cl a, cl b) { return !!sgn(det(a.t - a.s, b.s - a.s) - det(a.t - a.s, b.t - a.s)); }
bool turn_left(cp a, cp b, cp c) { return sgn(det(b - a, c - a)) > 0; }
bool point_on_segment(cp a,cl b) {
	return sgn (det(a - b.s, b.t - b.s)) == 0 && sgn (dot (b.s - a, b.t - a)) <= 0;}

int n;
vector <point> a;
map <point, int> id;
point at(int i) { return a[i % n]; }
int find_max (auto f) {
	int l = 0, r = (int) a.size() - 1;
	int d = !f(a[l], a[r]);
	if (d) swap(l, r);
	while (abs(r - l) > 1) {
		int mid = (l + r + d) / 2;
		if (f(a[mid], a[l]) && f(a[mid], a[mid + (d ? 1 : -1)])) l = mid;
		else r = mid; } return l;
}

vector <point> get_tan(point u) {
	if (id.count(u)) return {at(id[u] + n - 1), at(id[u] + 1)};
	if (point_on_segment(u, {a[0], a.back()})) {
		return {a.back(), a[0]};
	}
	return {
		at(find_max([&](cp x, cp y) {return turn_left(u, y, x);})),
		at(find_max([&](cp x, cp y) {return turn_left(u, x, y);}))};
}

void work() {
	a.clear(); id.clear();
	cin >> n;
	for (int i = 0; i < n; i++) {
		int x, y;
		cin >> x >> y;
		a.push_back({x, y});
	}
	for (int i = 0; i < n; i++) id[a[i]] = i;
	int Q;
	cin >> Q;
	while (Q--) {
		point u, v;
		u.read(); v.read();
		LD ans = 1e100;
 
		vector <point> uu, vv;
		uu = get_tan(u);
		vv = get_tan(v);
		if (det(uu[1] - u, v - u) > 0 && det(v - u, uu[0] - u) > 0
		&& det(vv[1] - v, u - v) > 0 && det(u - v, vv[0] - v) > 0
		) {
			for (auto i : uu) {
				for (auto j : vv) {
					if (has_inter({u, i}, {v, j})) {
						auto inter = line_inter({u, i}, {v, j});
						ans = min(ans, dis(u, inter) + dis(v, inter));
					}
				}
			}
		} else {
			ans = dis(u, v);
		}
		if (ans < 1e100) printf("%.9Lf\n", ans);
		else printf("-1\n");
	}
 
}
 
int main() {
	ios::sync_with_stdio(false); cin.tie(0);
	int T = 1;
//	cin >> T;
	for (int ca = 1; ca <= T; ca ++) {
		work();
	}
}