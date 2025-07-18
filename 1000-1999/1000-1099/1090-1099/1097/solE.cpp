#include <iostream>
#include <sstream>
#include <fstream>
#include <string>
#include <vector>
#include <deque>
#include <queue>
#include <stack>
#include <set>
#include <map>
#include <algorithm>
#include <functional>
#include <utility>
#include <bitset>
#include <cmath>
#include <cstdlib>
#include <ctime>
#include <cstdio>

using namespace std;

#define REP(i,n) for((i)=0;(i)<(int)(n);(i)++)
#define snuke(c,itr) for(__typeof((c).begin()) itr=(c).begin();itr!=(c).end();itr++)

#define INF (1<<29)
	
int N;
int a[100010];
bool covered[100010];

int x[100010];
int dp[100010];
int pre[100010];

int val_to_pos[100010];

vector <vector <int> > ANS;

int main3(int D, int pre_cut){
	int i;
	
	REP(i,N+1) x[i] = INF;
	x[0] = 0;
	
	REP(i,N) if(!covered[i]){
		int high = dp[i];
		int low = max(dp[i] - pre_cut, 1) - 1;
		
		while(high - low > 1){
			int mid = (low + high) / 2;
			if(x[mid] < a[i]){
				low = mid;
			} else {
				high = mid;
			}
		}
		
		dp[i] = high;
		pre[a[i]] = x[low];
		x[high] = a[i];
	}
	
	int maxdp = 0;
	int maxdppos = -1;
	REP(i,N) if(!covered[i]) if(dp[i] > maxdp){
		maxdp = dp[i];
		maxdppos = i;
	}
	
/*	REP(i,N) cout << a[i] << ' ';
	cout << endl;
	REP(i,N) cout << covered[i] << ' ';
	cout << endl;
	REP(i,N) cout << dp[i] << ' ';
	cout << endl;
*/	
	if(maxdp <= D){
		int sz = ANS.size();
		vector <int> empty;
	//	cerr << "maxdp: " << maxdp << endl;
		REP(i,maxdp) ANS.push_back(empty);
		REP(i,N) if(!covered[i]) ANS[sz + dp[i] - 1].push_back(a[i]);
		return -1;
	}
	
	int y = a[maxdppos];
	vector <int> v;
	
	while(y > 0){
		v.push_back(y);
		covered[val_to_pos[y]] = true;
		y = pre[y];
	}
	
	reverse(v.begin(),v.end());
	ANS.push_back(v);
	
	return v.size();
}

void main2(void){
	int i,j;
	
	REP(i,N) val_to_pos[a[i]] = i;
	
	ANS.clear();
	
	int D = 1;
	while((D + 1) * (D + 2) / 2 - 1 < N) D++;
	
	REP(i,N) covered[i] = false;
	REP(i,N) dp[i] = i + 1;
	
	int rem = N;
	int pre_cut = N + 5;
	
	while(1){
		int tmp = main3(D, pre_cut);
	//	cout << tmp << endl;
		if(tmp == -1) break;
		D--;
		rem -= tmp;
		if(rem == 0) break;
		pre_cut = tmp;
	}
	
	int sz = ANS.size();
	printf("%d\n", sz);
	
	REP(i,sz){
		int len = ANS[i].size();
		printf("%d", len);
		REP(j,len) printf(" %d", ANS[i][j]);
		printf("\n");
	}
}

int main(void){
	int Q,q,i;
	
	cin >> Q;
	REP(q,Q){
		scanf("%d", &N);
		REP(i,N) scanf("%d", &a[i]);
		main2();
	}
	/*
	while(1){
		N = 14;
		int i;
		REP(i,N) a[i] = i + 1;
		random_shuffle(a, a + N);
		
		REP(i,N) cout << a[i] << ' ';
		cout << endl;
	
		main2();
	} */
	
	return 0;
}