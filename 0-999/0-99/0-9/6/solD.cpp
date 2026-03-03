#include<bits/stdc++.h>
#define il inline
#define re register
#define ll long long
using namespace std;
#define MAXN 15
int n,A,B,a[MAXN],ans=1e9;
int b[MAXN],c[MAXN];

void DFS(int x,int sum){
	if(sum>=ans)return;
	if(x==n){
		// Fixed: Must check both archer n-1 and n
		if(a[n-1]<0 && a[n]<0){
			for(re int i=1;i<=n;++i)b[i]=c[i];
			ans=sum;
		}
		return;
	}
	// We start from the minimum shots needed to kill archer x-1
	int needed = 0;
	if(a[x-1] >= 0) {
		needed = a[x-1] / B + 1;
	}
	
	for(re int i=needed; i <= 16; ++i){
		a[x-1]-=i*B,a[x]-=i*A,a[x+1]-=i*B;
		// Archer x-1 is guaranteed to be dead here because of 'needed'
		c[x]=i;
		DFS(x+1,sum+i);
		a[x-1]+=i*B,a[x]+=i*A,a[x+1]+=i*B;
		
		// Optimization: if all three impacted archers are dead, 
		// taking more shots at 'x' is unlikely to be optimal unless x is the last possible shot position
		if (a[x-1] < 0 && a[x] < 0 && a[x+1] < 0 && x < n - 1) break;
	}
}

int main(){
	if(scanf("%d%d%d",&n,&A,&B) != 3) return 0;
	for(re int i=1;i<=n;++i) scanf("%d",&a[i]);
	
	// Pre-condition: Archer 1 and n MUST be killed by fireballs at 2 and n-1
	DFS(2,0);
	
	printf("%d\n",ans);
	bool first = true;
	for(re int i=1;i<=n;++i) {
		for(re int j=0;j<b[i];++j) {
			if(!first) printf(" ");
			printf("%d",i);
			first = false;
		}
	}
	printf("\n");
	return 0;
}
