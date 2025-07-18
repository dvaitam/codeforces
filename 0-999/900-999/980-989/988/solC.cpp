#include<iostream>
#include<cstdio>
#include<cstring>
#include<algorithm>
#include<cmath>
#define maxn 200005
#define LL long long
using namespace std;
int k,tot;
int len[maxn],sum[maxn],a[maxn],cnt;

struct check{
	int exc,fir,sec;
}c[maxn];

inline int rd(){
	int x=0,f=1; char c=' ';
	while(c>'9' || c<'0') {if(c=='-') f=-1;c=getchar();}
	while(c<='9' && c>='0') x=x*10+c-'0',c=getchar();
	return f*x; 
}

bool cmp(check a,check b){
	return a.exc<b.exc;
}

int main(){
	 k=rd();
	 for(int t=1;t<=k;t++){
	 	len[t]=rd();
	 	for(int i=1;i<=len[t];i++){
	 		a[i+cnt]=rd(); sum[t]+=a[i+cnt]; 
		 }
		for(int i=1;i<=len[t];i++){
			c[++tot].exc=sum[t]-a[i+cnt];
			c[tot].fir=t;
			c[tot].sec=i;
		}
		cnt+=len[t];
	 }
	 sort(c+1,c+tot+1,cmp);
	 for(int i=2;i<=tot;i++){
	 	if(c[i].exc==c[i-1].exc && c[i].fir!=c[i-1].fir){
	 		printf("YES\n");
	 		printf("%d %d\n",c[i-1].fir,c[i-1].sec);
	 		printf("%d %d\n",c[i].fir,c[i].sec);
	 		return 0;
		 }
	 }
	 printf("NO\n");
	 return 0;
}