#include<bits/stdc++.h>

using namespace std;



typedef long long ll;

ll t,n,k,seq[100009],ans[100009];

inline ll read(){

	ll s = 0,w = 1;

	char ch = getchar();

	while (ch > '9' || ch < '0'){ if (ch == '-') w = -1; ch = getchar();}

	while (ch <= '9' && ch >= '0') s = (s << 1) + (s << 3) + (ch ^ 48),ch = getchar();

	return s * w;

}



ll f(ll x,ll y,ll k){return (x - 1) * k + y;}

void get(ll n,ll k){

	if (!(k & 1)){

		for (ll i = 1;i <= n >> 1;i += 1) seq[(i << 1) - 1] = i,seq[i << 1] = n + 1 - i;

		return;

	}

	ll m = n / k,cur = 3 * m;

	for (ll i = 4;i <= k;i += 1){

		if (i & 1) for (ll j = m;j >= 1;j -= 1) seq[f(j,i,k)] = ++ cur;

		else for (ll j = 1;j <= m;j += 1) seq[f(j,i,k)] = ++ cur;

	}

	for (ll i = 1;i <= (m + 1 >> 1);i += 1){

		seq[f(i,1,k)] = (i << 1) - 1;

		seq[f(i,2,k)] = ((3 * m + 3) >> 1) - i;

		seq[f(i,3,k)] = 3 * m - i + 1;

	}

	for (ll i = (m + 3 >> 1);i <= m;i += 1){

		ll delta = i - (m + 3 >> 1);

		seq[f(i,1,k)] = ((3 * m + 3) >> 1) + delta;

		seq[f(i,2,k)] = (m << 1) + 1 + delta;

		seq[f(i,3,k)] = m - (m & 1) - (delta << 1);

	}

}



void print(){

	ll res = 0,sum = 0;

	for (ll i = 1;i <= k;i += 1) sum += ans[i];

	res = sum;

	for (ll i = k + 1;i <= n;i += 1) sum += ans[i] - ans[i - k],res = min(res,sum);

	printf("%lld\n",res);

	for (ll i = 1;i <= n;i += 1) printf("%lld ",ans[i]);

	puts("");

}

int main(){

	t = read();

	while (t --){

		n = read(),k = read();

		if (!(n % k)){

			get(n,k);

			for (ll i = 1;i <= n;i += 1) ans[i] = seq[i];

			print();

			continue;

		}

		ll q = n / k,r = n % k;

		if (r == 1){

			ll cur = 0,delta = (q << 1) + 1;

			for (ll i = 1;i <= n;i += k) ans[i] = ++ cur;

			for (ll i = n - k + 1;i >= 2;i -= k) ans[i] = ++ cur;

			get(q * (k - 2),k - 2),cur = 0;

			for (ll i = 3;i <= n;i += k) for (ll j = i;j <= i + k - 3;j += 1) ans[j] = seq[++ cur] + delta;

			print();

			continue;

		}

		if (k - r == 1){

			if (q == 1){

				ll cur = 0;

				ans[k] = n;

				get(n - 1,k - 1);

				for (ll i = 1;i < k;i += 1) ans[i] = seq[++ cur];

				for (ll i = k + 1;i <= n;i += 1) ans[i] = seq[++ cur];

				print();

				continue;

			}

			ll cur = n + 1,delta = q + 1;

			for (ll i = k;i <= n;i += k) ans[i] = -- cur;

			cur = 0;

			for (ll i = 1;i <= n;i += k) ans[i] = ++ cur;

			get((q + 1) * (r - 1),r - 1),cur = 0;

			for (ll i = 2;i <= n;i += k) for (ll j = i;j <= i + r - 2;j += 1) ans[j] = seq[++ cur] + delta;

			print();

			continue;

		}

		ll cur = 0,delta = (q + 1) * r;

		get((q + 1) * r,r);

		for (ll i = 1;i <= n;i += k) for (ll j = i;j <= i + r - 1;j += 1) ans[j] = seq[++ cur];

		get(q * (k - r),k - r),cur = 0;

		for (ll i = r + 1;i <= n;i += k) for (ll j = i;j <= i + (k - r) - 1;j += 1) ans[j] = seq[++ cur] + delta;

		print();

	}

	return 0;

}