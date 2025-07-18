#include <cstdio>

int buk[1 << 18 | 7], A[1 << 18 | 7];
int B[18], D[18], C;

inline void Ins(int x) {
	int y = x;
	for (int i = 17; ~i; --i) if (x >> i & 1) {
		if (!B[i]) { B[i] = x, D[i] = y, ++C; break; }
		x ^= B[i];
	}
}

void Solve(int *A, int j) {
	if (!j) return ;
	--j;
	A[1 << j] = D[j];
	Solve(A, j);
	Solve(A + (1 << j), j);
}

int main() {
	int N; scanf("%d", &N);
	while (N--) {
		int x; scanf("%d", &x);
		buk[x] = 1;
	}
	for (int i = 18; i >= 1; --i) {
		for (int j = 0; j < i; ++j) B[j] = 0;
		C = 0;
		for (int k = 0; k < 1 << i; ++k) if (buk[k]) Ins(k);
		if (C == i) {
			printf("%d\n", i);
			Solve(A, i);
			int x = 0;
			for (int k = 0; k < 1 << i; ++k) printf("%d ", x ^= A[k]);
			return 0;
		}
	}
	puts("0\n0 ");
	return 0;
}