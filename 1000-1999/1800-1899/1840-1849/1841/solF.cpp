#include<bits/stdc++.h>
using namespace std;
using ll=long long;
int n;
vector<pair<int, int> > v;

int main() {
	ios::sync_with_stdio(false);
	cin.tie(0);
	cin>>n;
	ll sx=0, sy=0;
	for(int a, b, c, d; n; --n) {
		cin>>a>>b>>c>>d;
		if((a==b)&&(c==d)) continue ;
		v.push_back({a-b, c-d});
		v.push_back({b-a, d-c});
		if(c<d||((c==d)&&(a<b))) {
			sx+=a-b; sy+=c-d;
		}
	}
	auto sec=[](const pair<int, int> &a)->int{
		if(a.first>0&&a.second>=0) return 1;
		else if(a.first<=0&&a.second>0) return 2;
		else if(a.first<0&&a.second<=0) return 3;
		else return 4; 
	};
	sort(v.begin(), v.end(), [&](const pair<int, int> &a, const pair<int, int> &b)->bool{
		int x=sec(a), y=sec(b);
		if(x!=y) return x<y;
		return 1ll*a.first*b.second>1ll*a.second*b.first;
	});
	double ans=1.0*sx*sx+1.0*sy*sy;
	for(auto t:v) {
		int x=t.first, y=t.second;
		sx+=x; sy+=y;
		ans=max(ans, 1.0*sx*sx+1.0*sy*sy);
	}
	cout<<fixed<<setprecision(10)<<ans<<"\n";
}