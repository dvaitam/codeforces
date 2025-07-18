#include<bits/stdc++.h>
using namespace std;
const int N=2e5+5;
int n,a[N],b[N],st[N],tp,m;
double ans[N],vb[N],va[N],o;
double calc(int x1,int x2,int y1,int y2){
	return double(n-1-x2-x2)*(y1-y2)/(x1-x2)+y2+y2;
}
int main(){
	scanf("%d",&n);
	for(int i=0;i<n;++i)scanf("%d",&a[i]);
	for(int i=0;i<n;++i)scanf("%d",&b[i]);
	m=n>>1,o=0;
	for(int i=n-1;~i;--i){
		while(tp>1&&1ll*(a[i]-a[st[tp-1]])*(st[tp]-i)<=1ll*(a[i]-a[st[tp]])*(st[tp-1]-i))--tp;
		ans[i]=max(ans[i],o);
		if(tp&&i<m&&st[tp]>=m)o=calc(i,st[tp],a[i],a[st[tp]]);
		st[++tp]=i;
		if(i==m){
			for(int j=0;j<m;++j){
				int l=1,r=tp;
				while(l<r){
					int mid=l+r+1>>1;
					if(1ll*(a[st[mid-1]]-b[j])*(st[mid]-j)>1ll*(a[st[mid]]-b[j])*(st[mid-1]-j))r=mid-1;
					else l=mid;
				}
				vb[j]=calc(j,st[l],b[j],a[st[l]]);
			}
		}
	}
	tp=0,o=0;
	for(int i=0;i<n;++i){
		while(tp>1&&1ll*(b[i]-b[st[tp-1]])*(i-st[tp])<=1ll*(b[i]-b[st[tp]])*(i-st[tp-1]))--tp;
		if(tp&&st[tp]<m&&i>=m)o=calc(st[tp],i,b[st[tp]],b[i]);
		ans[i]=max(ans[i],o);
		st[++tp]=i;
		if(i==m-1){
			for(int j=m;j<n;++j){
				int l=1,r=tp;
				while(l<r){
					int mid=l+r>>1;
					if(1ll*(a[j]-b[st[mid+1]])*(j-st[mid])<1ll*(a[j]-b[st[mid]])*(j-st[mid+1]))l=mid+1;
					else r=mid;
				}
				va[j]=calc(j,st[l],a[j],b[st[l]]);
			}
		}
	}
	for(int i=1;i<m;++i)vb[i]=max(vb[i],vb[i-1]);
	for(int i=n-2;i>=m;--i)va[i]=max(va[i],va[i+1]);
	for(int i=0;i<n;++i)printf("%.9lf%c",max(ans[i],(i<m)?vb[i]:va[i+1])," \n"[i==n-1]);
	return 0;
}