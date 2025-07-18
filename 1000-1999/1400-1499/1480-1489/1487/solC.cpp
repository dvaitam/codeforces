#include <cmath>
#include <stdio.h>
#include <cstring>
#include <iostream>
#include <algorithm>
using namespace std;

inline void file () {
	freopen("test.in", "r", stdin);
	freopen("test.out", "w", stdout);
}

inline int read () {
	int ans = 0, f = 1; char ch = getchar();
	while(!isdigit(ch)) { if(!(ch ^ '-')) f = -1; ch = getchar(); }
	while(isdigit(ch)) ans = ans * 10 + ch - '0', ch = getchar();
	return ans * f;
}

int n;
void solve () {
	n = read();
	if(n & 1) {
		for(int i = 1; i <= n; ++ i) {
			for(int j = i + 1; j <= n; ++ j) {
				if(j - (i + 1) + 1 <= (n >> 1)) putchar('1');
				else putchar('-'), putchar('1');
				putchar(' ');
			}
		}
	}
	else {
		for(int i = 1; i < n; i += 2) {
			printf("0 ");
			for(int j = i + 2; j <= n; ++ j) {
				if(j - (i + 2) + 1 <= (n - i - 1 >> 1)) putchar('1');
				else putchar('-'), putchar('1');
				putchar(' ');
			}
			for(int j = i + 2; j <= n; ++ j) {
				if(j - (i + 2) + 1 <= (n - i - 1 >> 1)) putchar('-'), putchar('1');
				else putchar('1');
				putchar(' ');
			}
		}
	}
	
	putchar('\n');
}

int main () {

	int T;
	T = read();
	while(T -- ) solve();

	return 0;
}