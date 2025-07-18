#include <bits/stdc++.h>

using namespace std;

vector<int> in(15);
int n;

bool dfs(int i, int angle){
	if(i == n){
		if(angle == 0) return true;
		else return false;
	}
	
	bool l,r;
	if(angle + in[i] >= 360) r = dfs(i+1, (angle+in[i])%360);
	else r = dfs(i+1, angle+in[i]);
	
	if(angle - in[i] < 0) l = dfs(i+1, 360 + angle-in[i]);
	else l = dfs(i+1, angle-in[i]);
	
	return l || r;
}

int main(){
	 // ios_base::sync_with_stdio(false);
	 // cin.tie(NULL);
	
	cin >> n;
	
	for(int i = 0; i < n; ++i)
		cin >> in[i];
	
	if(dfs(0, 0)) cout << "YES" << endl;
	else cout << "NO" << endl;
	
	return 0;
}