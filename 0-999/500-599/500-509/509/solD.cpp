#include<bits/stdc++.h> 

using namespace std;



const int maxn = 100 + 5;

int n, m;

int mat[maxn][maxn];

int a[maxn], b[maxn];



int gcd(int a, int b){

	return b == 0 ? a : gcd(b, a % b);

}



int main(){

	cin >> n >> m;

	for(int i = 1; i <= n; i++)

		for(int j = 1; j <= m; j++) scanf("%d", &mat[i][j]);

	for(int i = 1; i <= m; i++) b[i] = mat[1][i];

	for(int i = 1; i <= n; i++) a[i] = mat[i][1] - b[1];

	int g = 0;

	for(int i = 1; i <= n; i++)

		for(int j = 1; j <= m; j++){

			int x = abs(mat[i][j] - a[i] - b[j]);

			g = gcd(g, x);

		}

	if(!g) g = 1e9 + 1;

	for(int i = 1; i <= n; i++)

		for(int j = 1; j <= m; j++) if(g <= mat[i][j]) return 0*puts("NO");

	puts("YES");

	for(int i = 1; i <= n; i++) if(a[i] < 0) a[i] += g;

	printf("%d\n", g);

	for(int i = 1; i <= n; i++) printf("%d%c", a[i], i == n ? '\n' : ' ');

	for(int i = 1; i <= m; i++) printf("%d%c", b[i], i == m ? '\n' : ' ');

	return 0;

}