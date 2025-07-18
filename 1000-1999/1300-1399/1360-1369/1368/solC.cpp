using namespace std;

#define ll long long int

#define pb push_back

#define rb pop_back

#define ti tuple<int, int, int>

#define pii pair<int, int>

#define pli pair<ll, int>

#define pll pair<ll, ll>

#define mp make_pair

#define mt make_tuple

#include <iostream>

#include <vector>

#include <algorithm>

#include <queue>

#include <unordered_map>

#include <unordered_set>

#include <numeric>

#include <map>

#include <string>

#include <climits>

#include <cstring>

#include <cmath>

#include <iomanip> 

#include <set> 



vector<ll> subset;

string p;

vector<string> perms;

unordered_set<string> us;

void generate_subsets(vector<ll> v, ll cur, ll num) {

	if (cur == v.size()) {

		// process subset

		subset.push_back(num);

	}

	else {

		generate_subsets(v, cur + 1, num);

		num += v[cur];

		generate_subsets(v, cur + 1, num);

		num -= v[cur];

	}

}

vector<bool> visited(8, false);

void generate_permutations(string s) {

	if (p.size() == s.size()) {

		if (us.find(p) == us.end())

			perms.push_back(p);

		us.insert(p);

	}

	for (int i = 0; i < s.size(); i++) {

		if (!visited[i]) {

			p += (s[i]);

			visited[i] = true;

			generate_permutations(s);

			visited[i] = false;

			p.erase(p.size() - 1);

			//generate_permutations(s);

		}

	}

}

vector<int> primes;

unordered_set<int> pri;

void SieveOfEratosthenes(int n)

{

	// Create a boolean array "prime[0..n]" and initialize

	// all entries it as true. A value in prime[i] will

	// finally be false if i is Not a prime, else true.

	bool prime[1000001];

	memset(prime, true, sizeof(prime));



	for (int p = 2; p * p <= n; p++) {

		// If prime[p] is not changed, then it is a prime

		if (prime[p] == true) {

			// Update all multiples of p greater than or

			// equal to the square of it numbers which are

			// multiple of p and are less than p^2 are

			// already been marked.

			for (int i = p * p; i <= n; i += p)

				prime[i] = false;

		}

	}



	// Print all prime numbers

	for (int p = 2; p <= n; p++)

		if (prime[p]) {

			primes.push_back(p);

			pri.insert(p);

		}

}





vector<int> dive(int n)

{

	vector<int> div;

	// Note that this loop runs till square root

	for (int i = 1; i <= sqrt(n); i++)

	{

if (n % i == 0)

{

	// If divisors are equal, print only one

	if (n / i == i)

		div.push_back(i);



	else {

		div.push_back(i);

		div.push_back(n / i);

	}

}

	}

	return div;

}

unordered_map<ll, int> primeFactors(ll n)

{

	unordered_map<ll, int> m2;



	int count = 0;

	while (n % 2 == 0)

	{

		n = n / 2;

		count++;

	}

	if (count > 0)

		m2[2] = count;



	for (ll i = 3; i <= sqrt(n); i = i + 2)

	{

		count = 0;

		while (n % i == 0)

		{

			n = n / i;

			count++;

		}

		if (count > 0)

			m2[i] = count;



	}



	// This condition is to handle the case when n

	// is a prime number greater than 2

	if (n > 2) {



		m2[n] = 1;



	}

	return m2;

}

int main()

{

	//freopen("billboard.in", "r", stdin);

	//freopen("billboard.out", "w", stdout);

	ios_base::sync_with_stdio(false);

	cin.tie(NULL);

	int n;

	cin >> n;

	vector<pii> v;

	int x1 = 0, x2 = 1, y1 = 0, y2 = 1;

	cout << 4 + (3 * n) << '\n';

	cout << x1 << " " << y2 << '\n';

	cout << x1 << " " << y1 << '\n';

	cout << x2 << " " << y2 << '\n';

	cout << x2 << " " << y1 << '\n';

	int num = 4;

	while (n--) {

		num += 3;

		x2++;

		y2++;

		x1++;

		y1++;

		cout << x1 << " " << y2 << '\n';

		cout << x2 << " " << y2 << '\n';

		cout << x2 << " " << y1 << '\n';

	}



}