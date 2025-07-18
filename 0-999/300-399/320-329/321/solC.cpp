#include <iostream>
#include <string>
#include <vector>
#include <queue>
#include <algorithm>

using namespace std;
const int INF = 1000000000;

int main(){
	ios_base::sync_with_stdio(false);
	int n;
	cin >> n;
	vector< vector<int> > conn(n);
	vector<int> degrees(n);
	for(int i = 1; i < n; ++i){
		int a, b;
		cin >> a >> b;
		--a; --b;
		conn[a].push_back(b);
		conn[b].push_back(a);
		++degrees[a];
		++degrees[b];
	}
	vector<int> order;
	for(int i = 0; i < n; ++i){
		if(degrees[i] <= 1){ order.push_back(i); }
	}
	for(int i = 0; i < n; ++i){
		const int u = order[i];
		--degrees[u];
		for(int j = 0; j < conn[u].size(); ++j){
			if(--degrees[conn[u][j]] == 1){
				order.push_back(conn[u][j]);
			}
		}
	}
	vector<int> answer(n, -1);
	int maxval = 0;
	for(int i = 0; i < n; ++i){
		const int u = order[i];
		int count = 0, maximum = 0;
		for(int j = 0; j < conn[u].size(); ++j){
			const int v = conn[u][j];
			if(answer[v] < 0){ continue; }
			maximum = max(maximum, answer[v]);
			++count;
		}
		if(count == 0){
			answer[u] = 1;
		}else if(count == 1){
			answer[u] = maximum + 1;
		}else{
			answer[u] = (0x80000000u >> __builtin_clz(maximum)) << 1;
		}
		maxval = max(maxval, answer[u]);
	}
	if(maxval >= (1 << 26)){
		cout << "Impossible!" << endl;
	}else{
		for(int i = 0; i < n; ++i){
			if(i != 0){ cout << " "; }
			cout << static_cast<char>('Z' - __builtin_ctz(answer[i]));
		}
		cout << endl;
	}
	return 0;
}