#include <cstdio>

int k, n, p, A[31], S[901];

int main() {
	scanf("%d %d", &n, &k);
	for (int i = 1; i <= k; i++) {
		scanf("%d", A + i);
		S[A[i]] = 1;
	}
	p = 1;
	for (int i = 1; i <= k; i++) {
		printf("%d ", A[i]);
		for (int j = 1; j < n; p++) {
			if (!S[p]) {
				printf("%d ", p);
				j++;
			}
		}
	}
	puts("");
	return 0;
}