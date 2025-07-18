#include <bits/stdc++.h>
#define ZBJ TX
using namespace std;
using ll = long long;
const int inf = 0x3f3f3f3f;
const int maxn=305;
struct node {
	int l,r,id;
} tTT[maxn];
int n,m,a[maxn],b[maxn];
int main() {
	std::ios::sync_with_stdio(0);
	cin>>n>>m;
	int mx=-inf,mi=inf;
	for(int i=1; i<=n; ++i) {
		cin>>a[i];
		mx=max(mx,a[i]);
		mi=min(mi,a[i]);
	}
	for(int i=1; i<=m; ++i) {
		cin>>tTT[i].l>>tTT[i].r;
		tTT[i].id=i;
	}
	vector<int> vec;
	if(m==0 && n>=2) {
		cout<<mx-mi<<endl<<0<<endl;
	} else if(n==1) {
		cout<<0<<endl<<0<<endl;
	} else {
		int ans=-inf;
		vector<int> tp;
		for(int iii=1; iii<=n; ++iii) {
			tp.clear();
			memcpy(b,a,sizeof(a));
			for(int j=1; j<=m; ++j) {
				if(tTT[j].l<=iii&&tTT[j].r>=iii) {
					for(int k=tTT[j].l;k<=tTT[j].r;++k){
						b[k]--;
					}
					tp.push_back(tTT[j].id);
				}
			}
			mx=-inf,mi=inf;
			for(int i=1; i<=n; ++i) {
				mx=max(mx,b[i]);
				mi=min(mi,b[i]);
			}
			if(mx-mi>ans) {
				vec=tp;
				ans=mx-mi;
			}
		}
		cout<<ans<<endl;
		cout<<vec.size()<<endl;
		for(int i=0;i<vec.size();++i){
			cout<<vec[i]<< " ";
		}
		cout<<endl;

	}
	return 0;
}