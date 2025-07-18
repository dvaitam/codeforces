#include <cstdio>
#include <algorithm>
using namespace std;

int N, K, A, B;
void solve() {
	char small = 'G', big = 'B';
	if (A > B) swap(A, B), swap(small, big);
	// A <= B
	if ((B - 1) / K > A) {
		printf("NO\n");
		return;
	}
	int cnt = K;
	while (B > A) {
		if (!cnt)
			--A, printf("%c", small), cnt = K;
		else
			--B, printf("%c", big), --cnt;
	}
	while (A--) printf("%c%c", small, big);
	printf("\n");
}

int main() {
	scanf("%d%d%d%d", &N, &K, &A, &B);
	solve();
	return 0;
}