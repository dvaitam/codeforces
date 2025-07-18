#include<cstdio>
#include<iostream>
#include<algorithm>
#include<queue>

using namespace std;

const int N = 100005;
typedef long long ll;

int n,m,k;
ll f[N][2];
int md[N],mv[N];

struct env{
	int s,t,d,w,id;
	bool operator < (const env &x) const{
		return x.w>w || (x.w==w && x.d>d);
	}
}a[N];
bool cmp1(env x,env y) { return x.s<y.s; }
priority_queue<env> que;

int main()
{
	scanf("%d%d%d",&n,&m,&k);
	for(int i=1;i<=k;i++){
		scanf("%d%d%d%d",&a[i].s,&a[i].t,&a[i].d,&a[i].w);
		a[i].id=i;
	}
	sort(a+1,a+1+k,cmp1);
	
	int head=1;
	for(int i=1;i<=n;i++){
		while(head<=k && a[head].s<=i) que.push(a[head++]);
		while(!que.empty() && que.top().t<i) que.pop();
		if(que.empty()) {
			md[i]=i+1;
			mv[i]=0;
			continue;
		}
		md[i]=que.top().d+1;
		mv[i]=que.top().w;
	}
	
	for(int i=n;i>0;i--)
		f[i][0]=f[md[i]][0]+mv[i];
	for(int j=1;j<=m;j++)
		for(int i=n;i>0;i--)
			f[i][j&1]=min(f[i+1][(j-1)&1],f[md[i]][j&1]+mv[i]);
	printf("%lld\n",f[1][m&1]);
	
	return 0;
}