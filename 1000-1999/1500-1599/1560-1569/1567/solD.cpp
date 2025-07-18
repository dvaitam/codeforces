#include <iostream>

#include <map>

#include <set>

#include <queue>

#include <stack>

#include <algorithm>

#include <vector>

#include <string>

#include <iomanip>

#include <cmath>

#include <ctime>

#include <cstdio>

#include <cstdlib>

#include <cstring>

#include <climits>

//#include <unordered_map>

#define guo312 std::ios::sync_with_stdio(false), cin.tie(0), cout.tie(0)

#define ll long long

#define Inf LONG_LONG_MAX

#define inf INT_MAX

#define endl "\n"

#define PI 3.1415926535898

using namespace std;

int main(){

guo312;

	int t; cin>>t;

	while(t--){

		ll s,n; cin>>s>>n; ll base=1;

		while(base<=s){

			base*=10;

		}

		base/10;

		while(n){

			if(n==1){

				cout<<s<<endl;

				n--;

			}

			else if(s-base>=n-1){

				cout<<base<<" ";

				s-=base;

				n--;

			}

			else{

				base/=10;

			}

		}

	}

	return 0;

}