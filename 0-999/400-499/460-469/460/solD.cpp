#include <iostream>
#include <vector>

using namespace std;

typedef long long ll;

int main() {
	ll l, r;
	int k;
	cin >> l >> r >> k;
	ll result = l;
	vector<ll> S{l};
	for(int K = 2; K <= min(4, k); K++) {
		for(int i = 0; i < 2; i++) {
			if(l + i + K - 1 > r) continue;
			ll temp = 0;
			vector<ll> V;
			for(ll j = l + i; j < l + i + K; j++) {
				V.push_back(j);		
				temp ^= j;
			}
			if(temp < result) {
				result = temp;
				S = V;
			}
		}
	}
	if(k >= 3) {
		int msb = 62;
		for( ; msb >= 0; msb--) {
			if(l & (1ll << msb))
				break;
		}
		ll A = (1ll << (msb + 1)) | (1ll << msb); 
		ll B = (1ll << (msb + 1)) | (l ^ (1ll << msb));
		if(A <= r) {
			result = 0;
			S.clear();
			S.push_back(l);	
			S.push_back(A);
			S.push_back(B);
		}
	}

	cout << result << endl;
	cout << S.size() << endl;
	for(int i = 0; i < S.size(); i++)
		cout << S[i] << (i + 1 == S.size() ? "\n" : " ");
}