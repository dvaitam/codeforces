#include <cstdio>
#include <iostream>
#include <string>
#include <cstring>
#include <algorithm>
#include <vector>
using namespace std;

#define rep(i,a,n) for(int (i)=(a);(i)<(int)(n);(i)++)
#define foreach(c,itr) for(__typeof((c).begin()) itr=(c).begin();itr!=(c).end();itr++)
#define SZ(x) ((int) (x).size())

const int INF  =~0U >> 2;
const int N = 500 + 5;
int n,m;
char str[N][N], ans[N][N];
int dp[N][30][30],even[30],odd[30],cost[N][30][30];
pair<int,int> pre[N][30][30];

void calc(int row) {
	rep(i,0,26) even[i]=odd[i]=0;
	rep(i,0,m) rep(j,0,26) if(str[row][i]!='a'+j){
		if(i%2==0) even[j]++;
		else odd[j]++;
	}
	rep(i,0,26) rep(j,0,26) cost[row][i][j] = even[i] + odd[j];
}
int main(){
    //freopen("E.in","r",stdin);
	scanf("%d%d", &n,&m);
	rep(i,0,n) scanf("%s", str[i]);
	rep(i,0,n) calc(i);
	rep(i,0,n+1) rep(j,0,26) rep(k,0,26) dp[i][j][k] = INF;
	rep(i,0,26) rep(j,0,26) if(i!=j) dp[0][i][j] = 0;
	rep(i,0,n) {
		int a=0,b=1;
		rep(j,0,26) rep(k,0,26) if(dp[i][j][k]<dp[i][a][b]) a=j, b=k; 
		rep(j,0,26) rep(k,0,26) if(j!=k) {
			if(j!=a && k!=b) {
				dp[i+1][j][k] = dp[i][a][b] + cost[i][j][k];
				pre[i+1][j][k] = make_pair(a,b);
			} else {
				rep(s,0,26) rep(t,0,26) if(s!=j && t!=k) {
					if(dp[i][s][t]+cost[i][j][k] < dp[i+1][j][k]) {
						dp[i+1][j][k] = dp[i][s][t]+cost[i][j][k];
						pre[i+1][j][k] = make_pair(s,t);
					}
				}
			}
		}
	}
	int a=0, b=1;
	rep(i,0,26) rep(j,0,26) if(dp[n][i][j]<dp[n][a][b]) a=i, b=j;
	printf("%d\n",dp[n][a][b]);
	for(int i=n-1;i>=0;--i) {
		rep(j,0,m) if(j%2==0) ans[i][j]='a'+a; else ans[i][j]='a'+b;
		ans[i][m]='\0';
		int c = pre[i+1][a][b].first, d = pre[i+1][a][b].second;
		a = c, b = d;
	}
	rep(i,0,n) puts(ans[i]);
	
    return 0;
}