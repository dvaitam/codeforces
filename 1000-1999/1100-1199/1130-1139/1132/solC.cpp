#include<bits/stdc++.h>
#define ll long long 
#define lc o*2
#define rc o*2+1
using namespace std;
const int maxn = 5005;
int n,q;
struct node{
	int l,r;
	friend bool operator <(node a,node b){
		return a.l==b.l?a.r<b.r:a.l<b.l;
	}
}a[maxn];
int vis[5005];
int sum1[5005];
int sum2[5005];
int main(){
	cin>>n>>q;
	for(int i=1;i<=q;i++){
		scanf("%d%d",&a[i].l,&a[i].r);
	} 
	sort(a+1,a+1+q);
	for(int i=1;i<=q;i++){
		for(int j=a[i].l;j<=a[i].r;j++){
			vis[j]++;
		}
	}
	int sum=0;
	for(int i=1;i<=n;i++){
		if(vis[i])sum++;
		sum1[i]=sum1[i-1];
		sum2[i]=sum2[i-1]; 
		if(vis[i]==1)sum1[i]++;
		if(vis[i]==2)sum2[i]++; 
	}
	int maxx=0;
	for(int i=1;i<=q;i++){
		for(int j=i+1;j<=q;j++){
			int l1 = a[i].l;
			int l2 = a[j].l;
			int r1 = a[i].r;
			int r2 = a[j].r;
			int t=0;
			if(l2<=r1){
				t+=sum2[min(r1,r2)]-sum2[l2-1];
				t+=sum1[l2-1]-sum1[l1-1];
				t+=sum1[max(r2,r1)]-sum1[min(r1,r2)];
			}else{
				t+=sum1[r1]-sum1[l1-1];
				t+=sum1[r2]-sum1[l2-1];
			}
			maxx=max(sum-t,maxx); 
		}
	}
	cout<<maxx<<endl;
	return 0;
}