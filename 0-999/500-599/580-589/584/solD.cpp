/* Made by Ozhetov_A */           

#include <bits/stdc++.h>



#define F first

#define S second

#define mp make_pair

#define pb push_back

#define pf push_front

#define sz size()

#define all(x) x.begin(), x.end()



using namespace std;



typedef long long ll;

typedef long double ld;

typedef unsigned long long ull;



const ll N = 5e5;

const ll mod = 1e9 + 7;

const ll INF = 1e9 + 1;

const ld eps = 1e-9;

const ld pi = acos(-1.0);

    

ll n, a, b, c;



bool isprime(ll x) {

 for(ll i = 2; i * i <= x; ++i) {

  if(x % i == 0) {

   return false;

  }

 }

 return true;

}



int main () {

 ios_base::sync_with_stdio(false);

 cin.tie(0);

 cin >> n;

 for(ll i = n; i >= 2; --i) {

  if(isprime(i)) {

   a = i;

   break;

  }

 }

 for(ll i = 2; i <= n; ++i) {

  if(isprime(i)) {

   if(a + i != n && !isprime(n - (a + i))) {

    continue;

   } else

   if(a + i == n || isprime(n - (a + i))) {

    b = i;

    break;

   }

  }

 }

 if(a == n) {

  cout << 1 << "\n" << a;

 } else

 if(a + b == n) {

  cout << 2 << "\n" << a << " " << b;

 } else {

  cout << 3 << "\n" << a << " " << b << " " << n - (a + b);

 }

 return 0;

}