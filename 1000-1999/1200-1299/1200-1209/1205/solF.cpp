/* https://codeforces.com/blog/entry/78898 */
#include <stdio.h>

#define N	100
#define K	(N * (N + 1) / 2)

int dp[N + 1][K + 1];

void init() {
	int n, n_, k, k_, c;

	dp[1][1] = 1;
	for (n = 1; n <= N; n++)
		for (k = n; k <= n * (n + 1) / 2; k++)
			if (dp[n][k]) {
				for (c = 2; (n_ = n + c - 1) <= N; c++) {
					k_ = k + c * (c + 1) / 2 - 1;
					if (!dp[n_][k_])
						dp[n_][k_] = c;
				}
				for (c = 4; (n_ = n + c - 1) <= N; c++) {
					k_ = k + c;
					if (!dp[n_][k_])
						dp[n_][k_] = -c;
				}
			}
}

void trace(int *aa, int n, int k) {
	int n_, k_, c, i, j, a, a_, tmp;

	if (n == 1) {
		aa[0] = 1;
		return;
	}
	if (dp[n][k] > 0) {
		c = dp[n][k], n_ = n - c + 1, k_ = k - c * (c + 1) / 2 + 1;
		trace(aa, n_, k_);
		for (i = 0, j = n_ - 1; i < j; i++, j--)
			tmp = aa[i], aa[i] = aa[j], aa[j] = tmp;
		for (i = n_; i < n; i++)
			aa[i] = i + 1;
	} else {
		c = -dp[n][k], n_ = n - c + 1, k_ = k - c;
		trace(aa, n_, k_);
		for (i = 0; i < n_; i++)
			aa[i]++;
		if (c % 2 == 0) {
			for (a = n_ + 3; a <= n; a += 2)
				aa[i++] = a;
			aa[i++] = 1;
			for (a = n_ + 2; a <= n; a += 2)
				aa[i++] = a;
		} else {
			a_ = n_ + c / 2;
			for (a = n_ + 3; a < n; a += 2)
				aa[i++] = a + (a >= a_ ? 1 : 0);
			aa[i++] = a_, aa[i++] = 1;
			for (a = n_ + 2; a < n; a += 2)
				aa[i++] = a + (a >= a_ ? 1 : 0);
		}
	}
}

int main() {
	int q;

	init();
	scanf("%d", &q);
	while (q--) {
		static int aa[N];
		int n, k, i;

		scanf("%d%d", &n, &k);
		if (dp[n][k]) {
			trace(aa, n, k);
			printf("YES\n");
			for (i = 0; i < n; i++)
				printf("%d ", aa[i]);
			printf("\n");
		} else
			printf("NO\n");
	}
	return 0;
}