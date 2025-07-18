#include <bits/stdc++.h>
using namespace std;

using ii = pair<int, int>;
using ll = long long;

#define mp make_pair
#define mt make_tuple
#define pb push_back
#define eb emplace_back

int t, n, m, k, q;
char type;
bool grid[105][105];
vector<ii> ops;

void place_h(int r, int c) {
	assert(1 <= r && 1 <= n);
	assert(1 <= c && c + k - 1 <= m);
	for (int p = c; p <= c + k - 1; p++) {
		assert(!grid[r][p]);
		grid[r][p] = 1;
	}
}

void place_v(int r, int c) {
	assert(1 <= r && r + k - 1 <= n);
	assert(1 <= c && c <= m);
	for (int p = r; p <= r + k - 1; p++) {
		assert(!grid[p][c]);
		grid[p][c] = 1;
	}
}

void place(int r, int c, char type) {
	ops.eb(r, c);
	if (type == 'H') {
		place_h(r, c);
	} else {
		place_v(r, c);
	}
}

void clear_row(int r) {
	for (int c = 1; c <= m; c++) {
		assert(grid[r][c]);
		grid[r][c] = 0;
	}
}

void clear_col(int c) {
	for (int r = 1; r <= n; r++) {
		assert(grid[r][c]);
		grid[r][c] = 0;
	}
}

int main() {
	ios::sync_with_stdio(0);
	cin.tie(0);
	cin >> t;
	while (t--) {
		cin >> n >> m >> k >> q;
		if (k == 1) {
			int ptr = 1;
			for (int i = 1; i <= q; i++) {
				cin >> type;
				place(1, ptr, type);
				if (ptr == m) {
					clear_row(1);
					ptr = 1;
				} else {
					ptr++;
				}
			}
		} else if (k == min(n, m)) {
			if (k == max(n, m)) {
				for (int i = 1; i <= q; i++) {
					cin >> type;
					place(1, 1, type);
					if (type == 'H') {
						clear_row(1);
					} else {
						clear_col(1);
					}
				}
			} else if (k == n) {
				int ptr = 1;
				for (int i = 1; i <= q; i++) {
					cin >> type;
					if (type == 'H') {
						place(ptr, 2, type);
						if (ptr == n) {
							for (int j = 2; j <= 2 + k - 1; j++) {
								clear_col(j);
							}
							ptr = 1;
						} else {
							ptr++;
						}
					} else {
						place(1, 1, type);
						clear_col(1);
						for (int r = 1; r <= n; r++) {
							if (grid[r][2] && m == k + 1) {
								for (int c = 2; c <= m; c++) {
									grid[r][c] = 0;
									ptr = 1;
								}
							}
						}
					}
				}
			} else {
				int ptr = 1;
				for (int i = 1; i <= q; i++) {
					cin >> type;
					if (type == 'V') {
						place(2, ptr, type);
						if (ptr == m) {
							for (int j = 2; j <= 2 + k - 1; j++) {
								clear_row(j);
							}
							ptr = 1;
						} else {
							ptr++;
						}
					} else {
						place(1, 1, type);
						clear_row(1);
						for (int c = 1; c <= m; c++) {
							if (grid[2][c] && n == k + 1) {
								for (int r = 2; r <= n; r++) {
									grid[r][c] = 0;
									ptr = 1;
								}
							}
						}
					}
				}
			}
		} else {
			for (int i = 1; i <= q; i++) {
				cin >> type;
				bool h_clean = !(!grid[1][m - k + 1] && grid[1][m]);
				bool v_clean = !(!grid[n - k + 1][1] && grid[n][1]);
				if (type == 'H') {
					if (h_clean) {
						int to_place = -1;
						for (int i = 1; i <= n; i++) {
							if (!grid[i][m - k + 1]) {
								to_place = i;
								break;
							}
						}
						assert(to_place != -1);
						if (to_place < n - k) {
							place(to_place, m - k + 1, type);
						} else if (to_place == n - k) {
							place(to_place, m - k + 1, type);
							for (int j = m - k + 1; j <= m; j++) {
								if (grid[n][j]) {
									clear_col(j);
								}
							}
						} else if (v_clean) {
							place(to_place, m - k + 1, type);
							if (to_place == n - k + 1) {
								if (grid[n - k + 1][m - k]) {
									clear_row(n - k + 1);
								}
							} else {
								assert(!grid[n - k + 1][m - k]);
							}
							if (to_place == n) {
								for (int j = m - k + 1; j <= m; j++) {
									clear_col(j);
								}
							}
						} else {
							goto hell;
						}
					} else {
						hell:;
						int to_del = -1;
						for (int r = n - k + 1; r <= n; r++) {
							if (grid[r][m - k]) {
								to_del = r;
								break;
							}
						}
						if (to_del == -1) {
							int to_place = -1;
							for (int r = n - k + 1; r <= n; r++) {
								if (!grid[r][m - k + 1]) {
									to_place = r;
									break;
								}
							}
							assert(to_place != -1);
							place(to_place, m - k + 1, type);
							if (to_place == n) {
								for (int c = m - k + 1; c <= m; c++) {
									if (grid[1][c]) {
										clear_col(c);
									}
								}
							}
						} else {
							place(to_del, m - k + 1, type);
							clear_row(to_del);
						}
					}
				} else {
					if (v_clean) {
						int to_place = -1;
						for (int i = 1; i <= m; i++) {
							if (!grid[n - k + 1][i]) {
								to_place = i;
								break;
							}
						}
						assert(to_place != -1);
						if (to_place < m - k) {
							place(n - k + 1, to_place, type);
						} else if (to_place == m - k) {
							place(n - k + 1, to_place, type);
							for (int j = n - k + 1; j <= n; j++) {
								if (grid[j][m]) {
									clear_row(j);
								}
							}
						} else if (h_clean) {
							place(n - k + 1, to_place, type);
							if (to_place == m - k + 1) {
								if (grid[n - k][m - k + 1]) {
									clear_col(m - k + 1);
								}
							} else {
								assert(!grid[n - k][m - k + 1]);
							}
							if (to_place == m) {
								for (int j = n - k + 1; j <= n; j++) {
									clear_row(j);
								}
							}
						} else {
							goto hell2;
						}
					} else {
						hell2:;
						int to_del = -1;
						for (int c = m - k + 1; c <= m; c++) {
							if (grid[n - k][c]) {
								to_del = c;
								break;
							}
						}
						if (to_del == -1) {
							int to_place = -1;
							for (int c = m - k + 1; c <= m; c++) {
								if (!grid[n - k + 1][c]) {
									to_place = c;
									break;
								}
							}
							assert(to_place != -1);
							place(n - k + 1, to_place, type);
							if (to_place == m) {
								for (int r = n - k + 1; r <= n; r++) {
									if (grid[r][1]) {
										clear_row(r);
									}
								}
							}
						} else {
							place(n - k + 1, to_del, type);
							clear_col(to_del);
						}
					}
				}
			}
		}
		for (auto [r, c] : ops) {
			cout << r << ' ' << c << '\n';
		}
		ops.clear();
		for (int i = 1; i <= n; i++) {
			for (int j = 1; j <= m; j++) {
				grid[i][j] = 0;
			}
		}
	}
}