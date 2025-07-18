#include <stdio.h>
#include <math.h>

int prima[10000], cnt = 0, pertama, akar;

bool isprime(int a) {
	akar = sqrt(a);
	for (int j = 0; j < cnt; j++) {
		if (prima[j] > akar)
			break;
		if (a % prima[j] == 0) {
			pertama = j + 1;
			return false;
		}
	}
	return true;
}

int main() {
	int n;
	scanf("%d", &n);
	for (int i = 2; i <= n; i++) {
		if (i > 2)
			putchar(' ');
		if (isprime(i)) {
			prima[cnt++] = i;
			printf("%d", cnt);
		}
		else
			printf("%d", pertama);
	}
	printf("\n");
	return 0;
}