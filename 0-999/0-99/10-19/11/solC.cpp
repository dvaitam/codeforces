#include <cstdio>
#include <algorithm>
using namespace std;

const int nmax = 250;
int n, m;
char a[nmax + 2][nmax + 3];
int dr[] = {0, 1, 0, -1, 1, 1, -1, -1};
int dc[] = {1, 0, -1, 0, 1, -1, -1, 1};

int chk(int r, int c, int d) {
	int i, j;
	int d1 = (d + 1) % 4 + 4 * (d >= 4);
	int d2 = (d + 3) % 4 + 4 * (d >= 4);
	
	for (i = 1; ; i++) {
		r += dr[d], c += dc[d];
		if (a[r][c] == '0') return 0;
		if (d >= 4) {
			for (j = 0; j < 4; j++) if (a[r + dr[j]][c + dc[j]] == '1') return 0;
		}
//		printf("%d %d\n", r, c);
		if (a[r + dr[d2]][c + dc[d2]] == '1' || d < 4 && a[r + dr[d2] + dr[d]][c + dc[d2] + dc[d]] == '1') return 0;
		if (a[r + dr[d1]][c + dc[d1]] == '1') {
			// taisās beigties
			if (a[r + dr[d]][c + dc[d]] == '1') return 0;
			return i;
		} else {
			if (a[r + dr[d]][c + dc[d]] == '0') return 0;
		}
	}
}

int main() {
	int nt;
	
	for (scanf("%d", &nt); nt--; ) {
		int i, j, r = 0, t;
		
		scanf("%d%d", &n, &m);
		fill(a[0], a[0] + m + 2, '0');
		for (i = 1; i <= n; i++) {
			scanf("%s", a[i] + 1);
			a[i][0] = a[i][m + 1] = '0';
		}
		fill(a[n + 1], a[n + 1] + m + 2, '0');
		
		for (i = 1; i <= n; i++) for (j = 1; j <= m; j++) {
			if (a[i - 1][j - 1] == '0' && a[i - 1][j] == '0' && a[i - 1][j + 1] == '0'
			 && a[i][j - 1] == '0' && a[i][j] == '1' && a[i][j + 1] == '1'
			 && a[i + 1][j - 1] == '0' && a[i + 1][j] == '1'/* && a[i + 1][j + 1] == '0'*/) {
				r += (t = chk(i, j, 0)) && chk(i, j + t, 1) == t && chk(i + t, j + t, 2) == t && chk(i + t, j, 3) == t;
//				printf("%d+%d: %d %d %d %d\n", i, j, chk(i, j, 0), chk(i, j + t, 1), chk(i + t, j + t, 2), chk(i + t, j, 3));
			} else if (a[i - 1][j - 1] == '0' && a[i - 1][j] == '0' && a[i - 1][j + 1] == '0'
			 && a[i][j - 1] == '0' && a[i][j] == '1' && a[i][j + 1] == '0'
			 && a[i + 1][j - 1] == '1' && a[i + 1][j] == '0' && a[i + 1][j + 1] == '1') {
				r += (t = chk(i, j, 4)) && chk(i + t, j + t, 5) == t && chk(i + 2 * t, j, 6) == t && chk(i + t, j - t, 7) == t;
//				printf("%dx%d: %d %d %d %d\n", i, j, chk(i, j, 4), chk(i + t, j + t, 5), chk(i + 2 * t, j, 6), chk(i + t, j - t, 7));
			}
		}
		
		printf("%d\n", r);
	}
	
	return 0;
}