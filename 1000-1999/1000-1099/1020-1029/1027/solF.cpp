#include <bits/stdc++.h>
#define maxn 1000010
using namespace std;

inline int read(){
	int x=0,f=1;char cc=getchar();
	while(cc<'0' || cc>'9') {if(cc=='-') f=-1;cc=getchar();}
	while(cc>='0' && cc<='9') {x=x*10+cc-'0';cc=getchar();}
	return x*f;
}

int n,a[maxn],b[maxn],s1[maxn],s2[maxn],x[maxn],y[maxn],m,l,r;
int f[maxn<<1],v[maxn<<1],d[maxn<<1],c[maxn<<1];

inline bool cmp1(int i,int j){
	return a[i]<a[j];
}

inline bool cmp2(int i,int j){
	return b[i]<b[j];
}

inline int getfat(int k){
	if(f[k]==k) return k;
	f[k]=getfat(f[k]);
	return f[k];
}

int main(){
	n=read();
	for(int i=1;i<=n;i++){
		a[i]=read();b[i]=read();s1[i]=i;s2[i]=i;
	}
	sort(s1+1,s1+n+1,cmp1);
	sort(s2+1,s2+n+1,cmp2);
	l=0;r=0;
	m=0;
	while(l<n || r<n){
		int tmp=1000000007;
		if(l<n) tmp=min(tmp,a[s1[l+1]]);
		if(r<n) tmp=min(tmp,b[s2[r+1]]);
		m++;c[m]=tmp;
		while(l<n && a[s1[l+1]]==tmp){
			l++;x[s1[l]]=m;
		}
		while(r<n && b[s2[r+1]]==tmp){
			r++;y[s2[r]]=m;
		}
	}
	for(int i=1;i<=m;i++) f[i]=i;
	for(int i=1;i<=n;i++){
		if(getfat(x[i])==getfat(y[i])){
			d[f[x[i]]]++;
			continue;
		}
		d[f[y[i]]]+=d[f[x[i]]]+1;
		f[f[x[i]]]=f[y[i]];
	}
	for(int i=1;i<=m;i++) d[getfat(i)]--;
	for(int i=1;i<=m;i++) if(f[i]==i && d[i]>0){
		printf("-1");return 0;
	}
	for(int i=1;i<=m;i++) if(f[i]==i && d[i]==0) v[i]=1;
	for(int i=m;i;i--){
		if(!v[f[i]]){
			v[f[i]]=1;continue;
		}
		printf("%d",c[i]);return 0;
	}
}