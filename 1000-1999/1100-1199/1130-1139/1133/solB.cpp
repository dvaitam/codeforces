#include <bits/stdc++.h>
using namespace std;

inline void fastRead_int(int &x) {
    register int c = getchar();
    x = 0;
    int neg = 0;

    for(; ((c<48 || c>57) && c != '-'); c = getchar());

    if(c=='-') {
        neg = 1;
        c = getchar();
    }

    for(; c>47 && c<58 ; c = getchar()) {
        x = (x<<1) + (x<<3) + c - 48;
    }

    if(neg)
        x = -x;
}

int main(){
    
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);
	
	int n, k, temp;
	fastRead_int(n); fastRead_int(k);

	int buckets[k];
	memset(buckets, 0, sizeof buckets);

	for (int i = 0; i < n; ++i)	{
		fastRead_int(temp);
		buckets[ temp%k ]++;
	}

	long long int ans = 0;
	{
		for (int i = 1; i <= (k/2); ++i) {
			if(i == k-i) {
				ans += 2*(buckets[i] / 2);
			}
			else ans += 2*min(buckets[i], buckets[k-i]);
		}

		ans += 2*(buckets[0] / 2);
	}

	printf("%lld\n", ans);
    
	return 0;
}