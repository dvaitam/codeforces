#define _CRT_SECURE_NO_WARNINGS



#include <stdio.h>

#include <stdlib.h>

#include <string.h>

#include <iostream>

#include <vector>

#include <queue>

#include <string>

#include <algorithm>

#include <unordered_set>

#include <queue>

#include <unordered_map>

#include <set>

#include <map>

#include <stack>

#include <math.h>

#include <time.h>

#include <numeric>

#include <array>

#include <functional>

using namespace std;

typedef pair<int, int> pii;

typedef long long ll;

typedef pair<ll, ll> pll;

typedef array<ll, 3> arr3;





struct BIT {

	int n;

	vector<int> v;



	BIT(int _n) : n(_n), v(_n + 1, 0) {}



	int lowbit(int x) { return x & (-x); }



	void add(int i, int k) {

		while (i <= n) {

			v[i] += k;

			i += lowbit(i);

		}

	}



	int getSum(int i) {

		int res = 0;

		while (i > 0) {

			res += v[i];

			i -= lowbit(i);

		}

		return res;

	}



	int getSum(int l, int r) {

		return getSum(r) - getSum(l - 1);

	}

};

 

const ll base = 1e9 + 7;

 

ll pw(ll a, ll b) {

    ll t = a, r = 1;

    while (b) {

        if (b & 1) r = (r * t) % base;

        t = (t * t) % base;

        b >>= 1;

    }

    return r;

}

 

ll inv(ll x) {

    return pw(x, base - 2);

}

 

const int maxn = 1e6 + 42;

 

ll F[maxn], RF[maxn];

 

ll fact(ll n) {

    return F[n] ? F[n] : F[n] = n ? n * fact(n - 1) % base : 1;

}

ll rfact(ll n) {

    return RF[n] ? RF[n] : RF[n] = inv(fact(n));

}

 

ll nCr(ll n, ll r) {

    return fact(n) * rfact(r) % base * rfact(n - r) % base;

}



bool dp[41][41][41][41];

void solve() {

	ll T;

	cin >> T;

	const ll base = 1e9 + 7;

	for (ll ii = 1; ii <= T; ++ii) {

		int n, a, b;

		string s;

		cin >> n >> a >> b >> s;

		for (int i = 0; i <= n; ++i) {

			for (int j = 0; j < a; ++j) {

				for (int k = 0; k < b; ++k) {

					for (int l = 0; l <= n; ++l) {

						dp[i][j][k][l] = false;

					}

				}

			}

		}

		dp[0][0][0][0] = true;

		for (int i = 0; i < n; ++i) {

			int x = s[i] - '0';

			for (int j = 0; j < a; ++j) {

				for (int k = 0; k < b; ++k) {

					for (int l = 0; l <= i; ++l) {

						if (dp[i][j][k][l] == false) continue;

						dp[i + 1][(j * 10 + x) % a][k][l + 1] = true;

						dp[i + 1][j][(k * 10 + x) % b][l] = true;

					}

				}

			}

		}

		bool res = false;

		for (int i = 1; i < n; ++i) res |= dp[n][0][0][i];

		if (res) {

			int diff = 100;

			int red = -1;

			for (int l = 1; l < n; ++l) {

				if (dp[n][0][0][l] == false) continue;

				if (abs(l - (n - l)) < diff) {

					diff = abs(2 * l - n);

					red = l;

				}

			}

			string str;

			int state[4] = { n, 0, 0, red };

			for (int i = n - 1; i >= 0; --i) {

				int x = s[i] - '0';

				for (int j = 0; j < a; ++j) {

					if ((j * 10 + x) % a == state[1] && dp[i][j][state[2]][state[3] - 1]) {

						state[1] = j;

						state[3]--;

						str.push_back('R');

						goto CON;

					}

				}

				for (int k = 0; k < b; ++k) {

					if ((k * 10 + x) % b == state[2] && dp[i][state[1]][k][state[3]]) {

						state[2] = k;

						str.push_back('B');

						goto CON;

					}

				}

				CON:

				continue;

			}

			reverse(str.begin(), str.end());

			cout << str << endl;

		}

		else cout << "-1\n";

	}

}



int main() {

	solve();

    return 0;

}



//int main() {

//	int T = 1;

//	// cin >> T;

//	for (int ii = 0; ii < T; ++ii) {

//		int m, n;

//		cin >> m >> n;

//		vector<vector<int>> a(m, vector<int>(n));

//		for (int i = 0; i < m; ++i) {

//			for (int j = 0; j < n; ++j) {

//				cin >> a[i][j];

//			}

//		}

//		auto valid = [&](int x, int y) -> bool {

//			if (x < 0 || y < 0 || x >= m || y >= n) return false;

//			if (a[x][y] == 0) return false;

//			return true;

//		};

//		vector<vector<int>> res[5];

//		int d[][2] = { {0,0},{0,1},{1,0},{1,1} };

//		for (int i = 0; i < m - 1; ++i) {

//			for (int j = 0; j < n - 1; ++j) {

//				if (a[i][j] == 0) continue;

//				if (a[i][j] == a[i + 1][j] && a[i][j] == a[i][j + 1] && a[i][j] == a[i + 1][j + 1]) {

//					res[4].push_back(vector<int>({i, j, a[i][j]}));

//					a[i][j] = a[i + 1][j] = a[i][j + 1] = a[i + 1][j + 1] = 0;

//				}

//			}

//		}

//		if (res[4].size() == 0) {

//			printf("-1\n");

//			break;

//		}

//		auto dis = [&](const pii& p1, const pii& p2) -> int {

//			return min(abs(p1.first - p2.first), abs(p1.second - p2.second));

//		};

//		for (int i = 0; i < m - 1; ++i) {

//			for (int j = 0; j < n - 1; ++j) {

//				vector<pii> v;

//				if (valid(i, j) == false) continue;

//				v.push_back(pii(i, j));

//				if (valid(i + 1, j) && a[i][j] == a[i + 1][j]) v.push_back(pii(i + 1, j));

//				if (valid(i + 1, j + 1) && a[i][j] == a[i + 1][j + 1]) v.push_back(pii(i + 1, j + 1));

//				if (valid(i, j + 1) && a[i][j] == a[i][j + 1]) v.push_back(pii(i, j + 1));

//				if (v.size() == 3) {

//					res[3].push_back(vector<int>({i, j, a[i][j]}));

//					for (auto p : v) {

//						a[p.first][p.second] = 0;

//					}

//				}

//			}

//		}

//		for (int i = 1; i < m ; ++i) {

//			for (int j = 1; j < n; ++j) {

//				vector<pii> v;

//				if (valid(i, j) == false) continue;

//				v.push_back(pii(i, j));

//				if (valid(i - 1, j) && a[i][j] == a[i - 1][j]) v.push_back(pii(i - 1, j));

//				if (valid(i, j - 1) && a[i][j] == a[i][j - 1]) v.push_back(pii(i, j - 1));

//				if (v.size() == 3) {

//					res[3].push_back(vector<int>({i, j, a[i][j]}));

//					for (auto p : v) {

//						a[p.first][p.second] = 0;

//					}

//				}

//			}

//		}

//		for (int i = 0; i < m - 1; ++i) {

//			for (int j = 0; j < n - 1; ++j) {

//				vector<pii> v;

//				if (valid(i, j) == false) continue;

//				v.push_back(pii(i, j));

//				if (valid(i + 1, j) && a[i][j] == a[i + 1][j]) v.push_back(pii(i + 1, j));

//				if (valid(i, j + 1) && a[i][j] == a[i][j + 1]) v.push_back(pii(i, j + 1));

//				if (v.size() == 2) {

//					res[2].push_back(vector<int>({i, j, a[i][j]}));

//					for (auto p : v) {

//						a[p.first][p.second] = 0;

//					}

//				}

//			}

//		}

//		for (int i = 0; i < m - 1; ++i) {

//			for (int j = 0; j < n - 1; ++j) {

//				vector<pii> v;

//				if (valid(i, j) == false) continue;

//				res[1].push_back(vector<int>({i, j, a[i][j]}));

//			}

//		}

//		for (auto &it: res) reverse(it.begin(), it.end());

//		for (auto it : res[1]) printf("%d %d %d\n", it[0] + 1, it[1] + 1, it[2]);

//		for (auto it : res[2]) printf("%d %d %d\n", it[0] + 1, it[1] + 1, it[2]);

//		for (auto it : res[3]) printf("%d %d %d\n", it[0] + 1, it[1] + 1, it[2]);

//		for (auto it : res[4]) printf("%d %d %d\n", it[0] + 1, it[1] + 1, it[2]);

//		

//	}

//	return 0;

//}



/**

struct BIT {

	int n;

	vector<int> v;



	BIT(int _n) : n(_n), v(_n + 1, 0) {}



	int lowbit(int x) { return x & (-x); }



	void add(int i, int k) {

		while (i <= n) {

			v[i] += k;

			i += lowbit(i);

		}

	}



	int getSum(int i) {

		int res = 0;

		while (i > 0) {

			res += v[i];

			i -= lowbit(i);

		}

		return res;

	}



	int getSum(int l, int r) {

		return getSum(r) - getSum(l - 1);

	}

};



class Solution {

public:

	int countRangeSum(vector<int> a, int lower, int upper) {

		int n = a.size();

		vector<ll> pre(n + 1);

		pre[0] = 0;

		for (int i = 0; i < n; ++i) pre[i + 1] = pre[i] + (ll)a[i];

		int res = 0;

		map<ll, int> mp;

		for (auto it : pre) mp[it] = 0;

		int ap = 1;

		for (auto& it : mp) it.second = ap++;



		BIT bt(1 << 17);

		bt.add(mp[0], 1);

		ll maxx = prev(mp.end())->first;

		ll minn = mp.begin()->first;

		for (int i = 1; i <= n; ++i) {

			ll up = pre[i] - lower;

			ll lo = pre[i] - upper;

			

			if (!(up < minn || lo > maxx)) { //  注意处理区间包含情况

				int l = mp.lower_bound(lo)->second;

				int r = prev(mp.upper_bound(up))->second;

				res += bt.getSum(l, r);

			}

			bt.add(mp[pre[i]], 1);

		}

		return res;

	}

};



int main() {

	Solution ss;

	cout << ss.countRangeSum(vector<int>({ 1,1,0,-2}), -3, 1);

	return 0;

}

*/



/**

int main() {

	int T;

	cin >> T;

	for (int ii = 1; ii <= T; ++ii) {

		int n, k;

		cin >> n >> k;

		int mask = n - 1;

		if (n == 4 && k == 3 || n == 2 && k == 1) printf("-1\n");

		else {

			vector<pii> res;

			vector<bool> vis(n, false);

			if (k != n - 1) {

				res.push_back(pii(n - 1, k));

				vis[n - 1] = vis[k] = true;

				for (int i = n - 1; i >= 0; --i) {

					if (vis[i]) continue;

					int other = (~i) & mask;

					if (vis[other]) other = 0;

					res.push_back(pii(i, other));

					vis[i] = vis[other] = true;

				}

			}

			else {

				res.push_back(pii(n - 1, n - 2));

				res.push_back(pii(1, n - 3));

				vis[n - 1] = vis[n - 2] = vis[1] = vis[n - 3] = true;

				for (int i = n - 4; i >= 0; --i) {

					if (vis[i]) continue;

					int other = (~i) & mask;

					if (vis[other]) other = 0;

					res.push_back(pii(i, other));

					vis[i] = vis[other] = true;

				}

			}

			for (int i = 0; i < res.size(); ++i) {

				printf("%d %d\n", res[i].first, res[i].second);

			}

		}

	}

	return 0;

}

*/

/**

int d[][2] = { {0,1}, {0,-1},  {1,0}, {-1,0} };



struct Record {

	int x, y, bx, by;



	Record() {}

	Record(int a, int b, int c, int d) :

		x(a), y(b), bx(c), by(d) {}



	Record push(int index, bool& b, vector<vector<char>>& grid) {

		int dx = d[index][0];

		int dy = d[index][1];

		b = true;

		Record res = *this;

		res.x += dx;

		res.y += dy;

		if (grid[res.x][res.y] == '#') b = false;

		if (res.x == res.bx && res.y == res.by) {

			res.bx += dx;

			res.by += dy;

			if (grid[res.bx][res.by] == '#') b = false;

		}

		//b = true;

		return res;

	}



	bool end(vector<vector<char>>& grid) {

		return grid[bx][by] == 'T';

	}



	long long hash() {

		long long c1 = x * 32 + y;

		long long c2 = bx * 32 + by;

		return c1 * 1024 + c2;

	}

};



class Solution {

public:

	int minPushBox(vector<vector<char>>& g) {

		vector<vector<char>> grid(g.size() + 2, vector<char>(g[0].size() + 2, '#'));

		for (int i = 0; i < g.size(); ++i) {

			for (int j = 0; j < g[0].size(); ++j) {

				grid[i + 1][j + 1] = g[i][j];

			}

		}

		int m = grid.size();

		int n = grid[0].size();

		Record p0;

		unordered_set<long long> st;

		for (int i = 0; i < m; ++i) {

			for (int j = 0; j < n; ++j) {

				switch (grid[i][j]) {

				case '.':

				case 'T':

				case '#': break;

				case 'B':

					p0.bx = i;

					p0.by = j;

					grid[i][j] = '.';

					break;

				case 'S':

					p0.x = i;

					p0.y = j;

					grid[i][j] = '.';

					break;

				}

			}

		}

		queue<Record> q[2];

		int now = 0;

		q[now].push(p0);

		int cnt = 0;

		while (q[now].size()) {

			while (q[now].size()) {

				Record ap = q[now].front();

				q[now].pop();

				long long code = ap.hash();

				if (st.find(code) != st.end()) continue;

				st.insert(code);

				if (ap.end(grid)) return cnt;

				for (int i = 0; i < 4; ++i) {

					bool b = true;

					Record r = ap.push(i, b, grid);

					if (b) {

						if (r.bx == ap.bx && r.by == ap.by) q[now].push(r);

						else q[!now].push(r);

					}

				}

			}

			cnt++;

			now = !now;

		}

		return -1;

	}

};



int main() {

	Solution ss;

	string g[][8] = {

		{"#",".",".","#","#","#","#","#" },

		{"#",".",".","T","#",".",".","#" },

		{"#",".",".",".","#","B",".","#" },

		{"#",".",".",".",".",".",".","#" },

		{"#",".",".",".","#",".","S","#" },

		{"#",".",".","#","#","#","#","#" }

	};

	int m = sizeof(g) / sizeof(g[0]);

	int n = sizeof(g[0]) / sizeof(g[0][0]);

	vector<vector<char>> grid(m, vector<char>(n));

	for (int i = 0; i < m; ++i) {

		for (int j = 0; j < n; ++j) {

			grid[i][j] = g[i][j][0];

		}

	}

	cout << ss.minPushBox(grid);

	return 0;

}

*/

/*

vector<ll> mi;



ll diff(ll x) {

	ll up = lower_bound(mi.begin(), mi.end(), x) - mi.begin();

	return mi[up] - x;

}



int main() {



	for (ll i = 0; i < 62; ++i) {

		mi.push_back(1ll << i);

	}



	int T;

	cin >> T;

	for (int ii = 1; ii <= T; ++ii) {

		ll n;

		cin >> n;

		vector<ll> input(n);

		for (int i = 0; i < n; ++i) scanf("%d", &input[i]);

		sort(input.begin(), input.end());

		vector<ll> v;

		v.push_back(1);

		for (int i = 1; i < n; ++i) {

			if (input[i] == input[i - 1]) v.back()++;

			else v.push_back(1);

		}

		ll res = diff(n) + 2;		//  one group

		int m = v.size();



		ll sum = 0;

		for (int i = 0; i < m; ++i) {

			sum += v[i];

			res = min(res, diff(sum) + diff(n - sum) + 1);		// two groups

		}



		sum = 0;





		cout << res << endl;

	}

	return 0;

}

*/

/**

class Solution {

public:

	bool canDistribute(vector<int>& nums, vector<int>& a) {

		unordered_map<int, int> mp;

		for (auto it : nums) mp[it]++;

		vector<int> v;

		for (auto& it : mp) v.push_back(it.second);

		sort(v.begin(), v.end(), greater<int>());

		sort(a.begin(), a.end(), greater<int>());

		while (v.size() > a.size()) v.pop_back();



		int n = a.size();

		int m = v.size();

		

		

	}

};

*/



/**



class Solution {

public:

	int minFlips(vector<vector<int>>& mat) {

		int m = mat.size();

		int n = mat[0].size();

		queue<vector<vector<int>>> q[2];

		int now = 0;

		unordered_set<int> visit;

		auto hash = [&](vector<vector<int>>& matrix) -> int {

			int res = 0;

			for (int i = 0; i < m; ++i) {

				for (int j = 0; j < n; ++j) {

					int index = (i << 3) + j;

					res |= matrix[i][j] * (1 << index);

				}

			}

			return res;

		};

		auto change = [&](vector<vector<int>>& matrix, int x, int y) -> void {

			matrix[x][y] = !matrix[x][y];

			if (x > 0) matrix[x - 1][y] = !matrix[x - 1][y];

			if (x < m - 1) matrix[x + 1][y] = !matrix[x + 1][y];

			if (y > 0) matrix[x][y - 1] = !matrix[x][y - 1];

			if (y < n - 1) matrix[x][y + 1] = !matrix[x][y + 1];

		};

		int res = 0;

		q[now].push(mat);

		visit.insert(hash(mat));

		while (q[now].size()) {

			while (q[now].size()) {

				auto ap = q[now].front();

				q[now].pop();

				if (hash(ap) == 0) return res;

				for (int i = 0; i < m; ++i) {

					for (int j = 0; j < n; ++j) {

						vector<vector<int>> matrix = ap;

						change(matrix, i, j);

						int code = hash(matrix);

						if (visit.find(code) == visit.end()) {

							visit.insert(code);

							q[!now].push(matrix);

						}

					}

				}

			}

			res++;

			now = !now;

		}

		return -1;

	}

};



int main() {

	Solution ss;

	vector<int> nums({ 1,1,2,2,3,3,4,4,5,5,6,6,7,7,8,8,9,9,10,10,11,11,12,12,13,13,14,14,15,15,16,16,17,17,18,18,19,19,20,20,21,21,22,22,23,23,24,24,25,25,26,26,27,27,28,28,29,29,30,30,31,31,32,32,33,33,34,34,35,35,36,36,37,37,38,38,39,39,40,40,41,41,42,42,43,43,44,44,45,45,46,46,47,47,48,48,49,49,50,50});

	vector<int> a({ 2,2,2,2,2,2,2,2,2,3 });

	//cout << ss.canDistribute(nums, a) << endl;

	int in[][2] = { {0, 0},{0, 1} };

	int m = sizeof(in) / sizeof(in[0]);

	int n = sizeof(in[0]) / sizeof(in[0][0]);

	vector<vector<int>> mat(m, vector<int>(n));

	for (int i = 0; i < m; ++i) {

		for (int j = 0; j < n; ++j) {

			mat[i][j] = in[i][j];

		}

	}

	cout << ss.minFlips(mat) << endl;

	return 0;

}

*/