#include<bits/stdc++.h>

using namespace std;

const int N = 53;

int n, m;

char str[N][N];

void solve(){

	cin >> n >> m;

	bool flg[2] = {1, 1};

	for(int i = 0;i < n;++ i){

		cin >> str[i];

		for(int j = 0;j < m;++ j)

			if(str[i][j] == 'R') flg[(i ^ j) & 1] = 0;

			else if(str[i][j] == 'W') flg[!((i ^ j) & 1)] = 0;

	}

	if(flg[0]){

		puts("YES");

		for(int i = 0;i < n;++ i){

			for(int j = 0;j < m;++ j)

				putchar(((i ^ j) & 1) ? 'R' : 'W');

			putchar('\n');

		}

	} else if(flg[1]){

		puts("YES");

		for(int i = 0;i < n;++ i){

			for(int j = 0;j < m;++ j)

				putchar(((i ^ j) & 1) ? 'W' : 'R');

			putchar('\n');

		}

	} else puts("NO");

}

int main(){

	ios::sync_with_stdio(0);

	int T; cin >> T; while(T --) solve();

}