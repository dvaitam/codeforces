#include<cstdio>
#include<iostream>
#include<cstring>
#include<algorithm>
#include<cmath>
#include<queue>
using namespace std;
typedef long long ll;
int getint(){
	int x=0,f=1; char ch=getchar();
	while(ch>'9'||ch<'0') {if(ch=='-')f=-1;ch=getchar();}
	while(ch<='9'&&ch>='0') {x=x*10+ch-'0';ch=getchar();}
	return f*x;
} 
#define N 1123456
struct qujian{
	int l,r,id;
}a[N];
bool cmp(qujian a,qujian b) {
	return a.l<b.l;
}
bool cmp2(qujian a,qujian b){
	return a.r>b.r;
}
priority_queue <int ,vector<int > ,greater<int > >  q;
int tot=0;
int r[N];
qujian aa[N];
int n;
int main(){
	n=getint();int k=getint();
	for(int i=1;i<=n;i++) a[i].l=getint(),a[i].r=getint(),a[i].id=i;
	sort(a+1,a+1+n,cmp); 
//	for(int i=1;i<=n;i++) cout<<a[i].l<<" "<<a[i].r<<endl;
	for(int i=1;i<=n;i++){
		if(tot==k){
			int t=q.top(); 
			if(t<=a[i].r) {
				q.pop(); q.push(a[i].r);
	//			cout<<"pp "<<t<<endl; 
	//			cout<<"p "<<a[i].r<<endl;
			}
			r[i]=q.top();
		}
		else {
			
			q.push(a[i].r); ++tot;
			if(tot!=k)r[i]=-1234567890;
			else r[i]=q.top(); 
	//		cout<<"p "<<a[i].r<<endl; 
		} 
	}
	//for(int i=1;i<=n;i++) cout<<r[i]<<" "; cout<<endl;
	int ans=0,ansi=1;
	for(int i=1;i<=n;i++) {
		int t=max(0,r[i]-a[i].l+1); if(t>ans) { ans=t;ansi=i;}  
	}
	int tot=0;
	for(int i=1;i<=ansi;i++) aa[++tot]=a[i];
	sort(aa+1,aa+1+tot,cmp2);
	printf("%d\n",ans);
	if(ans==0){
		for(int i=1;i<=k;i++) printf("%d ",i);
		 return 0;
	}
	for(int i=1;i<=k;i++) printf("%d ",aa[i].id);
	return 0;
}