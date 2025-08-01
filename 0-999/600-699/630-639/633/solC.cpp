#include<bits/stdc++.h>
using namespace std;

char s[10010], a[1111111], *word[100010];
int n, m, len, u[1111111][26], cnt = 1, vis[1111111], dp[10010];
vector<int> ans;

int calc(int i) {
	if(i == n) return 1;
	if(dp[i] != -1) return dp[i];
	int x = 1;
	for(int idx, j=i ; j<n ; ++j) {
		idx = s[j] - 'a';
		if(!u[x][idx]) return dp[i] = 0;
		x = u[x][idx];
		if(vis[x]){
			if(calc(j+1) == 1) {
				ans.push_back(vis[x]);
				return dp[i] = 1;
			}
		}
	}
	return dp[i] = 0;
}

int main(int argc, char **argv) {
#ifndef ONLINE_JUDGE
	freopen("a.in", "r", stdin);
#endif
	scanf("%d%s%d", &n, s, &m);
	reverse(s, s+n);
	for(int l=0, x, i=1 ; i<=m ; ++i) {
		scanf("%s", a+l);
		word[i] = (a+l);
		len = strlen(a+l);
		l += len+1;
		x = 1;
		for(int idx, j=0 ; j<len ; ++j) {
			idx = (word[i][j] | 32) - 'a';
			if(u[x][idx] == 0)
				u[x][idx] = ++cnt;
			x = u[x][idx];
		}
		vis[x] = i;
	}
	memset(dp, -1, sizeof dp);
	calc(0);
	for(int i=0 ; i<(int)ans.size() ; ++i)
		printf("%s ", word[ans[i]]);
	return 0;
}