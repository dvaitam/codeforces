#include <bits/stdc++.h>
using namespace std;
int t,n,a[1005],b[1005];
vector<pair<int,int> > ans;
void Print(){
	printf("%d\n",ans.size());
	for(auto p:ans) printf("%d %d\n",p.first,p.second);
}
void Move(int l,int r,int x){
	if(x<0) x+=r-l+1;
	for(int i=l;i<=r;++i) b[(i+x-l)%(r-l+1)+l]=a[i];
	for(int i=l;i<=r;++i) a[i]=b[i];
}
signed main(){
	scanf("%d",&t);
//	int cnt=0;
	while(t--){
		ans.clear();
		scanf("%d",&n);
		for(int i=1;i<=n;++i) cin>>a[i];
//		++cnt;
//		if(cnt==218){
//			cout<<"!!!!! n="<<n<<endl;
//			for(int i=1;i<=n;++i) cout<<a[i]<<" ";
//			cout<<endl;
//			break;
//		}
//		if(cnt>3) continue;
		bool flag=1;
		for(int i=1;i<=n;++i)
			if(a[i]!=i) flag=0;
		if(flag){
			printf("%d\n0\n",n);
			continue;
		}
		flag=1;
		for(int i=1;i<n;++i)
			if(a[i+1]-a[i]==1||a[i+1]-a[i]==1-n);
			else flag=0;
		if(flag){
			printf("%d\n",n-1);
			int p;
			for(int i=1;i<=n;++i)
				if(a[i]==1) p=i;
			printf("%d\n",p-1);
			for(int i=1;i<p;++i) printf("2 1\n");
			continue;
		}
		int inv=0;
		for(int i=1;i<n;++i)
			for(int j=i+1;j<=n;++j)
				if(a[i]>a[j]) ++inv;
		if((inv&1)&&((n-1)&1)){
			printf("%d\n",n-3);
			int posn;
			for(int i=1;i<=n;++i)
				if(a[i]==n) posn=i;
			if(posn==1){
				ans.push_back({2,1});
				Move(1,n-2,-1);
				posn=n-2;
			}else if(posn==2){
				ans.push_back({1,2});
				Move(1,n-2,1);
				posn=3;
			}
			for(int i=posn;i<n;++i) ans.push_back({3,4});
			Move(3,n,n-posn);
			--n;
//			for(int i=1;i<=n;++i) cout<<a[i]<<" ";
//			cout<<endl;
			if(a[1]==1){
				ans.push_back({1,2});
				int tmp=a[n-1];
				for(int i=n-1;i>=2;--i) a[i]=a[i-1];
				a[1]=tmp;
			}else if(a[n]==1){
				ans.push_back({2,3});
				for(int i=n;i>=3;--i) a[i]=a[i-1];
				a[2]=1;
			}
			for(int i=2;i<=n-2;++i){
				int pos1,posi;
				for(int j=1;j<=n;++j)
					if(a[j]==1) pos1=j;
				for(int j=1;j<=n;++j)
					if(a[j]==i) posi=j;
				if(posi==1){
					if(pos1==2){
						ans.push_back({1,2});
						int tmp=a[n-1];
						for(int j=n-1;j>=2;--j) a[j]=a[j-1];
						a[1]=tmp;
						++pos1;++posi;
					}else{
						ans.push_back({2,1});
						for(int j=1;j<=n-2;++j) a[j]=a[j+1];
						a[n-1]=i;
						--pos1;
						posi=n-1;
					}
				}
				if(posi!=n){
					if(posi<pos1){
						for(int j=2;j<=posi;++j) ans.push_back({3,2});
						Move(2,n,-(posi-1));
						pos1-=posi-1;
						posi=n;
					}else{
						for(int j=posi;j<n;++j) ans.push_back({2,3});
						Move(2,n,n-posi);
						pos1+=n-posi;
						posi=n;
					}
				}
				for(int j=pos1+i-2;j<n-1;++j) ans.push_back({1,2});
				Move(1,n-1,n-1-(pos1+i-2));
				pos1+=n-1-(pos1+i-2);
				ans.push_back({3,2});
				Move(2,n,-1);
			}
			ans.push_back({2,1});
			Move(1,n-1,-1);
			if(a[n-1]==n-1){
				Print();
				continue;
			}else{
				for(int i=1;i<=(n-3)/2;++i){
					ans.push_back({1,2});
					ans.push_back({2,3});
					ans.push_back({2,3});
					ans.push_back({2,1});
					ans.push_back({3,2});
					ans.push_back({2,1});
					ans.push_back({2,3});
					ans.push_back({2,3});
				}
				ans.push_back({2,3});
				Print();
			}
		}else{
			printf("%d\n",n-2);
			if(a[1]==1){
				ans.push_back({1,2});
				int tmp=a[n-1];
				for(int i=n-1;i>=2;--i) a[i]=a[i-1];
				a[1]=tmp;
			}else if(a[n]==1){
				ans.push_back({2,3});
				for(int i=n;i>=3;--i) a[i]=a[i-1];
				a[2]=1;
			}
			for(int i=2;i<=n-2;++i){
				int pos1,posi;
				for(int j=1;j<=n;++j)
					if(a[j]==1) pos1=j;
				for(int j=1;j<=n;++j)
					if(a[j]==i) posi=j;
				if(posi==1){
					if(pos1==2){
						ans.push_back({1,2});
						int tmp=a[n-1];
						for(int j=n-1;j>=2;--j) a[j]=a[j-1];
						a[1]=tmp;
						++pos1;++posi;
					}else{
						ans.push_back({2,1});
						for(int j=1;j<=n-2;++j) a[j]=a[j+1];
						a[n-1]=i;
						--pos1;
						posi=n-1;
					}
				}
				if(posi!=n){
					if(posi<pos1){
						for(int j=2;j<=posi;++j) ans.push_back({3,2});
						Move(2,n,-(posi-1));
						pos1-=posi-1;
						posi=n;
					}else{
						for(int j=posi;j<n;++j) ans.push_back({2,3});
						Move(2,n,n-posi);
						pos1+=n-posi;
						posi=n;
					}
				}
				for(int j=pos1+i-2;j<n-1;++j) ans.push_back({1,2});
				Move(1,n-1,n-1-(pos1+i-2));
				pos1+=n-1-(pos1+i-2);
				ans.push_back({3,2});
				Move(2,n,-1);
			}
			ans.push_back({2,1});
			Move(1,n-1,-1);
			if(a[n-1]==n-1){
				Print();
				continue;
			}else{
				for(int i=1;i<=(n-3)/2;++i){
					ans.push_back({1,2});
					ans.push_back({2,3});
					ans.push_back({2,3});
					ans.push_back({2,1});
					ans.push_back({3,2});
					ans.push_back({2,1});
					ans.push_back({2,3});
					ans.push_back({2,3});
				}
				ans.push_back({2,3});
				Print();
			}
		}
	}
	return 0;
}