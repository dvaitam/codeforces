#include <cstdio>

#include <algorithm>

#include <vector>

using namespace std;

const int N=300000;

int st[N],en[N],n,x,y,size[N];

double ans=0,sst,sen;

vector<int> v[N];

void dfs(int x,int fx){

	size[x]=1;

	for(int i=0;i<v[x].size();i++){

	    if(v[x][i]!=fx){

	    	dfs(v[x][i],x);

	    	size[x]+=size[v[x][i]];

	    	st[x]+=st[v[x][i]];

	    	ans+=1.0*st[v[x][i]]*size[v[x][i]]*en[x];

	    }

	}ans+=(sst-st[x])*(n-size[x])*en[x];

} 

int main(){

	scanf("%d",&n);

	for(int i=1;i<n;i++){

		scanf("%d%d",&x,&y);

		v[x].push_back(y);

		v[y].push_back(x);

	}for(int i=1;i<=n;i++){

		scanf("%d%d",st+i,en+i);

		sst+=st[i]; sen+=en[i];

	}dfs(1,-1);

	printf("%.11f\n",ans/sst/sen);

	return 0;

}