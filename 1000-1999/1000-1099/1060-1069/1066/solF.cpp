#include<cstring>
#include<iostream>
#include<cstdio>
#include<cmath>
#include<cstring>
#include<string>
#include<algorithm>
#include<vector>
#include<queue>
#include<stack>
#include<map>
#include<set>
using namespace std;
#define fi first
#define se second
#define mp make_pair
#define pb push_back
#define rep(i,a,b) for (int i=(a);i<(b);i++)
#define per(i,a,b) for (int i=(b)-1;i>=a;i--)
#define sz(a) (int) a.size()
#define de(a) cout<<#a<<" = "<<a<<endl
#define dd(a) cout<<#a<<" = "<<a<<" "
typedef long long ll;
typedef pair<int,int> pii;
typedef vector<int> vi;
//-----
const double eps = 1e-8;
const int inf = 0x7f7f7f7f;
const int N=2e5+5;
const int mod=998244353;
ll dp[N];
struct P{
	int x,y;
}p[N];
int le(int d) {
	return max(p[d].x,p[d].y);
}
int get(int a,int b) {
	return abs(p[a].x - p[b].x) + abs(p[a].y - p[b].y);
}
bool operator < (P a,P b) {
	int f1 = max(a.x,a.y),f2 = max(b.x,b.y);
	if(f1 != f2) return f1 < f2;
	if(a.x != b.x) return a.x < b.x;
	return a.y > b.y;
}
int main()
{
	int n;scanf("%d",&n);
	rep(i,1,n+1)
	scanf("%d%d",&p[i].x,&p[i].y);
	sort(p+1,p+n+1);
	//rep(i,1,n+1) cout<<p[i].x<<' '<<p[i].y<<endl;
	memset(dp,0x3f,sizeof(dp));
	int pl = 0,pr = 0;
	dp[0] = 0;
	while(pr < n) {
		int d = pr+1;
		while(d < n && le(d+1) == le(pr+1)) d++;
		dp[d] = min(dp[pl] + get(pl,pr+1),dp[pr] + get(pr,pr+1)) + get(pr+1,d);
		dp[pr+1] = min(dp[pl] + get(pl,d),dp[pr] + get(pr,d)) + get(pr+1,d);
		pl = pr + 1;
		pr = d;
	//	dd(pl),de(pr);
	//	dd(dp[pl]),de(dp[pr]);
	}
	printf("%I64d\n",min(dp[pl],dp[pr]));
	return 0;
}