#include <cmath>
#include <ctime>
#include <cstdio>
#include <cctype>
#include <cstdlib>
#include <cstring>
#include <iostream>
#include <algorithm>
#include <set>
#include <map>
#include <stack>
#include <queue>
#include <string>
#include <vector>
#define maxl 1000000000
#define maxn 201000
using namespace std;

int first[maxn],nxt[maxn];
pair<int,int> mm[maxn],ans[maxn];


void tj(int x,int y){
	nxt[x]=first[y];
	first[y]=x;
}

int main(){
	int n,s,t1,i,x,j,k,mid;
	scanf("%d%d",&n,&s);
	if(s&1){
		printf("No\n");
		return 0;
	}
	for(i=1;i<=n;++i){
		scanf("%d",&x);
		tj(i,x);
	}
	t1=0;
	for(i=s;i>=1;--i)while(first[i]!=0){
		x=first[i];
		first[i]=nxt[x];
		j=i;
		for(k=1;k<=i;++k){
			while(j>0 && first[j]==0)--j;
			if(j==0){
				printf("No\n");
				return 0;
			}
			mid=first[j];
			mm[k]=make_pair(mid,j);
			first[j]=nxt[mid];
			ans[++t1]=make_pair(x,mid);
		}
		for(k=1;k<=i;++k)tj(mm[k].first,mm[k].second-1);
	}
	printf("Yes\n%d\n",t1);
	for(i=1;i<=t1;++i)printf("%d %d\n",ans[i].first,ans[i].second);
	return 0;
}