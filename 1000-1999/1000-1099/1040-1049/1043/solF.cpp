#include<bits/stdc++.h>
using namespace std;
#define SET(a,e) memset(a,e,sizeof(a))
#define LL long long
#define LD long double
#define pb push_back
#define x first
#define y second
#define PII pair<int,int>
#define PLI pair<LL,int>
#define PIL pair<int,LL>
#define PLL pair<LL,LL>
#define PDD pair<LD,LD>
#define eps 1e-9
#define HH1 402653189
#define HH2 1610612741

int n, a[300055], sieve[300055];
bool have[10][300055];
vector<int> arr, ini;
vector<int> pos[20];

void build(int x) {
	arr.clear();
	for (int i = 1; i <= 300000; i++)
		if (have[x][i])
			arr.push_back(i);
}

int main() {
/*	
	sieve[1] = 1;
	for (int i = 2; i <= 300000; i++) {
		if (sieve[i])
			continue;
		for (int j = i + i; j <= 300000; j += i)
			sieve[j] = i;
	}
*/	
	scanf("%d", &n);
	
	for (int i = 0; i < n; i++) {
		scanf("%d", a + i);
		if (a[i] == 1) {
			printf("1\n");
			return 0;
		}
		have[0][a[i]] = true;
	}
	
	int g = a[0];
	
	for (int i = 1; i < n; i++)
		g = __gcd(g, a[i]);
	
	if (g > 1) {
		printf("-1\n");
		return 0;
	}
	
	for (int i = 2; i <= 300000; i++) {
		if (!have[0][i])
			continue;
		for (int j = i + i; j <= 300000; j += i)
			have[0][j] = false;
	}
	
	build(0);
	for (int i = 1; i <= 300000; i++)
		if (have[0][i]) {
			ini.push_back(i);
//			printf("i = %d\n", i);
		}
	
	for (int _ = 1; _ < 10; _++) {
		for (int i : arr) {
			for (int j : ini) {
				int g = __gcd(i, j);
				if (g == 1) {
					printf("%d\n", _ + 1);
					return 0;
				}
				have[_][g] = true;
			}
		}
		
		build(_);
	}
	
	return 255;
	
	return 0;
}