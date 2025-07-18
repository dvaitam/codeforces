#include<ctime>
#include<iostream>
#include<algorithm>
#include<cstdio>
#include<cstdlib>
#include<cmath> 
#include<cstring> 
#include<cassert>
#include<string>
#include<sstream>
#include<fstream>
#include<deque>
#include<queue>
#include<vector>
#include<map>
#include<list>
#include<stack>
#include<set>
#include<bitset>
#include<iomanip>
#include<utility>
#include<functional>
#include<cctype>
#include<cerrno>
#include<cfloat>
#include<ciso646>
#include<climits>
#include<clocale>
#include<complex>
#include<csetjmp>
#include<csignal>
#include<cstdarg>
#include<cstddef>
#include<cwchar>
#include<cwctype>
#include<exception>
#include<locale>
#include<numeric>
#include<new>
#include<stdexcept>
#include<limits>
using namespace std;

#define ll long long
#define INF 1e9
#define rep(i,n) for(int (i)=0;(i)<n;i++)
#define REP(i,n) for(int (i)=1;(i)<=n;i++)
#define mk(a,b) make_pair(a,b)
#define fi first
#define se second
#define pii pair<int,int>
#define sz(s) s.size()
#define all(s) (s.begin(),s.end())

const int maxn=2505;
int n;
class point{
public:
	int x,y;
	double operator -(const point &other)const{
		return sqrt(double(abs((double)x-(double)other.x)*abs((double)x-(double)other.x))+(double)(abs((double)y-(double)other.y)*abs((double)y-(double)other.y)));
	}
}p[maxn];
double dp[maxn][maxn][2];
double ans;

int main(){
	scanf("%d",&n);
	rep(i,n)scanf("%d%d",&p[i].x,&p[i].y);
	REP(len,n-1)rep(i,n){
		int end=i-len;
		if(end<0)end+=n;
		int start=i+1;
		if(start>=n)start-=n;
		dp[len+1][start][0]=max(dp[len+1][start][0],dp[len][i][0]+(p[start]-p[i]));
		dp[len+1][end][1]=max(dp[len+1][end][1],dp[len][i][0]+(p[end]-p[i]));
		end=i+len;
		if(end>=n)end-=n;
		start=i-1;
		if(start<0)start+=n;
		dp[len+1][start][1]=max(dp[len+1][start][1],dp[len][i][1]+(p[start]-p[i]));
		dp[len+1][end][0]=max(dp[len+1][end][0],dp[len][i][1]+(p[end]-p[i]));
	}
	rep(i,n)rep(k,2)ans=max(ans,dp[n][i][k]);
	printf("%.15f",ans);
	return 0;
}