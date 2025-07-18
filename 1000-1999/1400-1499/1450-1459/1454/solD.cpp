#include <bits/stdc++.h>
using namespace std;

using ll = long long;
using ld = long double;
using ull = unsigned long long;

#define alpha ios::sync_with_stdio(0); cin.tie(0);

int MAX = 1e9 + 7;
// ull fact[1001];

bool comp(int a, int b);
ull power(ull x, int y, int p);
ull modInverse(ull n, int p);
// ull nCrModP(ull n, int r, int p);
bool isPrime(ull n);
ull gcd(ull a, ull b);

pair<ull, ull> primeFac(ull n) {
	pair<ull, ull> temp;
	ull p = 0;
	ull c = 0;
	ull prime_num = 0;
	ull count = 0;

	prime_num = 2;
	while (n % 2 == 0)  
    {  
        count++; 
        n = n/2;  
    }

    if(count > c) {
    	c = count;
    	p = prime_num;
    }

    for(ull i=3; i<= sqrt(n); i=i+2) {
    	prime_num = i;
		count = 0;

		while (n % i == 0)  
	    {  
	        count++; 
	        n = n/i;  
	    }

	    if(count > c) {
	    	c = count;
	    	p = prime_num;
	    }
    }

    temp.first = c;
    temp.second = p;

    return temp;
}

void solve() {
	ull n; cin >> n;

	if(isPrime(n)) {
		cout << 1 << '\n';
		cout << n;
		return;
	}
	
	pair<ull, ull> v = primeFac(n);
	ull prime_num = v.second;
	ull count = v.first;

	cout << count << '\n';

	for(ull i=1; i<= count - 1; i++) {
		cout << prime_num << " ";
		n = n/prime_num;
	}

	cout << n;

	
}




int main(int argc, char const *argv[])
{	
	alpha;
	#ifndef ONLINE_JUDGE 
	freopen("input.txt", "r", stdin); 
	freopen("error.txt", "w", stderr);
	freopen("output.txt", "w", stdout); 
	#endif

	int t=1;
	cin >> t;

	// to precompute factorial upto certain upper bound
	// fact[0] = 1;
	// for (int i = 1; i <= 1000; i++) fact[i] = (fact[i - 1] * i) % MAX;

	while(t--) {
		solve();
		cout << '\n';
	}
	cerr<<"time taken : "<<(float)clock()/CLOCKS_PER_SEC<<" secs"<<endl;
	return 0;
}






// -------------------------- AUX Functions ------------------------------ //
bool comp(int a, int b) {
	return a < b;
}

// to compute (x^y) % p
ull power(ull x, int y, int p) {
	ull res = 1;
	x = x % p;
	while(y>0) {if(y & 1) res = (res * x) % p;
		y = y >> 1;	// y=y/2;
		x = (x * x) % p;}
	return res;
}

// to compute (n^(-1)) % p
ull modInverse(ull n, int p) {return power(n, p-2, p);}

// ull nCrModP(ull n, int r, int p) {
// 	if(r==0) return 1;
//  return (fact[n] * modInverse(fact[r], p) % p * modInverse(fact[n-r], p) % p) % p;
// }

bool isPrime(ull n) {
	if(n<=1) return false;
	if(n<=3) return true;
	if(n%2 == 0 || n%3 == 0) return false;

	for(ull i=5; i*i <= n; i=i+6) {if(n % i == 0 || n % (i+2) == 0) return false;}
	return true;
}

ull gcd(ull a, ull b) {
	while(a%b != 0) {
		ull temp = a;
		a = b;
		b = temp % b;
	}
	return b;
}