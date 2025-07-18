#include <cassert>
#include <cctype>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <algorithm>
#include <bitset>
#include <iomanip>
#include <iostream>
#include <map>
#include <queue>
#include <set>
#include <sstream>
#include <string>
#include <utility>
#include <vector>
typedef long long LL;
typedef unsigned long long ULL;
typedef long double LD;
using std::cin; using std::cout;
using std::endl;
using std::bitset; using std::map;
using std::queue; using std::priority_queue;
using std::set; using std::string;
using std::stringstream; using std::vector;
using std::pair; using std::make_pair;
typedef pair<int, int> pii;
typedef pair<LL, LL> pll;
typedef pair<ULL, ULL> puu;
#ifdef DEBUG
using std::cerr;
#define pass cerr << "[" << __FUNCTION__ << "] : line = " << __LINE__ << endl;
#define display(x) cerr << #x << " = " << x << endl;
#define displaya(a, st, n)                      \
	{                                           \
		cerr << #a << " = {";                   \
		for (int qwq = (st); qwq <= (n); ++qwq) \
			if (qwq == (st))                    \
				cerr << a[qwq];                 \
			else                                \
				cerr << ", " << a[qwq];         \
		cerr << "}" << endl;                    \
	}
#define displayv(a) displaya(a, 0, (int)(a.size() - 1))
#define eprintf(...) fprintf(stderr, __VA_ARGS__)
#include <ctime>
class MyTimer {
	clock_t st;
public:
	MyTimer() { cerr << std::fixed << std::setprecision(0); reset(); }
	~MyTimer() { report(); }
	void reset() { st = clock_t(); }
	void report() {  cerr << "Time consumed: " << (clock() - st) * \
		1e3 / CLOCKS_PER_SEC << "ms" << endl; }
} myTimer;
#else
#define cerr if(false) std::cout
#define pass ;
#define display(x) ;
#define displaya(a, st, n) ;
#define displayv(a) ;
#define eprintf(...) if(0) fprintf(stderr, __VA_ARGS__)
class MyTimer {
public: void reset() {} void report() {}
} myTimer;
#endif

template<typename A, typename B>
std::ostream& operator << (std::ostream &cout, const pair<A, B> &x) {
	return cout << "(" << x.first << ", " << x.second << ")";
}
template<typename T1, typename T2>
inline bool chmin(T1 &a, const T2 &b) { return a > b ? a = b, true : false; }
template<typename T1, typename T2>
inline bool chmax(T1 &a, const T2 &b) { return a < b ? a = b, true : false; }

const int maxN = 100000 + 233;
int n, m;
int a[maxN], b[maxN], c[maxN], d[maxN], e[maxN];
int ans[maxN];

struct Q {
	int t;
	int x, y, L, R;
	Q(int t, int x, int y, int L, int R) :
		t(t), x(x), y(y), L(L), R(R) {}
	bool operator < (const Q &rhs) const {
		return x != rhs.x
			? x < rhs.x
			: y != rhs.y ? y < rhs.y : t > rhs.t;
	}
};
vector<Q> q, tmp;

void initialize() {
	std::ios::sync_with_stdio(false);
	cin >> n >> m;
	for(int i = 1; i <= n; ++i) cin >> a[i];
	for(int i = 1; i <= n; ++i) cin >> b[i];
	for(int i = 1; i <= n; ++i) cin >> c[i];
	for(int i = 1; i <= m; ++i) cin >> d[i];
	for(int i = 1; i <= m; ++i) cin >> e[i];
	for(int i = 1; i <= n; ++i) {
		q.emplace_back(i, -c[i] + a[i], c[i] + a[i], a[i], b[i]);
	}
	for(int i = 1; i <= m; ++i) {
		q.emplace_back(-i, -e[i] + d[i], d[i] + e[i], d[i], d[i]);
	}
	vector<int> t;
	for(auto &p : q) t.push_back(p.L), t.push_back(p.R);
	std::sort(t.begin(), t.end());
	t.erase(std::unique(t.begin(), t.end()), t.end());
	for(auto &p : q)
		p.L = std::lower_bound(t.begin(), t.end(), p.L) - t.begin() + 1,
		p.R = std::lower_bound(t.begin(), t.end(), p.R) - t.begin() + 1;
	std::sort(q.begin(), q.end());
	tmp.resize(q.size(), Q(0, 0, 0, 0, 0));
}

int C[maxN * 3];
void add(int p, int x) { p += 2;
	for(int i = p; i <= n * 2 + m + 5; i += i & -i) C[i] += x;
}
int sum(int p) { p += 2;
	int r = 0;
	for(int i = p; i > 0; i -= i & -i) r += C[i];
	return r;
}
void solve(int L, int R) {
	if(L >= R) return;
	int M = (L + R) >> 1;
	solve(L, M); solve(M + 1, R);
	int i = L, j = M + 1, qwq = L;
	const auto left = [&](const Q &p) {
		tmp[qwq++] = p;
		if(p.t > 0) add(p.L, 1), add(p.R + 1, -1);
	};
	const auto right = [&](const Q &p) {
		tmp[qwq++] = p;
		if(p.t < 0) ans[-p.t] += sum(p.L);
	};
	while(i <= M && j <= R) {
		if(q[i].y < q[j].y || (q[i].y == q[j].y && q[i].t > 0))
			left(q[i++]);
		else
			right(q[j++]);
	}
	while(i <= M) left(q[i++]);
	while(j <= R) right(q[j++]);
	for(int k = L; k <= M; ++k) if(q[k].t > 0) add(q[k].L, -1), add(q[k].R + 1, 1);
	for(int k = L; k <= R; ++k) q[k] = tmp[k];
}

int main() {
	initialize();
	solve(0, (int)q.size() - 1);
	for(int i = 1; i <= m; ++i) printf("%d%c", ans[i], " \n"[i == m]);
	return 0;
}