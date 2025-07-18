#include <bits/stdc++.h>
using namespace std;

#define vi vector < int >
#define vll vector < ll >
#define pb push_back
#define ll long long int
#define llu unsigned long long
#define MOD 1000000007
#define F first
#define S second
#define INF 2000000000
#define dbg(x) { cout<< #x << ": " << (x) << endl; }
#define loop(i,a,b) for(int i = a ; i < b ;++i)
#define all(x) x.begin(),x.end()
#define MAX 1000007
#define mapsi map < string, int >
#define mapll map < ll, ll >

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);
    ll n;
    ll x,y;
    cin >> n >> x >> y;
    if(x == n && y == n){
        cout << "Black";
    }
    else if(max(abs(x-1),abs(y-1)) > max(abs(n-x),abs(n-y))){
        cout << "Black";
    }
    else{
        cout << "White";
    }
	return 0;
}