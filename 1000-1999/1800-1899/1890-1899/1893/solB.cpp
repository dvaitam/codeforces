// LUOGU_RID: 156267999
// LUOGU_RID: 154919549

#include <bits/stdc++.h>

using ll = long long;
const ll mod = 998244353;
//#define int long long

//#define std::pair<int,int> pii



void solve() {
	int n,m;std::cin>>n>>m;
	std::vector<int> a(n + 2),b(m + 2),lis(n + 2),ha(n + 2);
	for(int i = 1;i <= n;i++)std::cin>>a[i],ha[i] = a[i];
	for(int i = 1;i <= m;i++)std::cin>>b[i];
	std::sort(b.begin() + 1,b.begin() + m + 1, std::greater<int>());
	std::vector<int> bit(n + 2);	
	std::sort(ha.begin() + 1,ha.begin() + n + 1);
	
	std::map<int,int > mp;
	
	int cnt = 0;
	
	for(int i = 1;i <= n;i++){
		if(!mp[ha[i]]) mp[ha[i]] = ++cnt;
	}
	
	auto lowbit = [&](int x){
		return x & -x;
	};
	
	auto add = [&](int i,int x){
		while(i <= n){
			bit[i] = std::min(x,bit[i]);
			i += lowbit(i);
		}
	};
	
	auto get_min = [&](int i){
		int res = 0x7ffffff;
		while(i){
			res = std::min(res,bit[i]);
			i -= lowbit(i);
		}
		return res;
	};
	
	/*for(int i = 1;i <= n;i++){
		lis[i] = get_min(mp[a[i]] - 1) + 1;
		add(mp[a[i]],lis[i]); 
	}*/
	
	int p = 1;
	
	std::vector<int> ans;
	
	for(int i = 1;i <= m;i++){
		while(p <= n && a[p] >= b[i]) {
			ans.push_back(a[p]);
			p++;
		}
		ans.push_back(b[i]);
	}
	
	//ans.push_back(a[n]);
	if(p <= n) for(int i = p;i <= n;i++)ans.push_back(a[i]);
	
	for(auto i : ans)std::cout<<i<<" ";
	std::cout<<'\n';
	return;
}

signed main() {
    std::ios::sync_with_stdio(false);
    std::cin.tie(0);
    int t = 1;
    std::cin >> t;
    //  init();
    while (t--) {
        solve();
        // std::cout.flush();
    }
    // std::cout<<5 * inv(2ll) % mod<<'\n';
    return 0;
}