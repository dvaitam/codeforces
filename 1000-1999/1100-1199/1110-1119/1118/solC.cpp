#include<bits/stdc++.h>
using namespace std;

#define int long long int
#define double long double
#define pb push_back
#define pii pair<int,int>
#define fi first
#define se second
#define rep(i,a,b) for (int i=a; i<b; ++i)
#define dbg(x) { cerr<<#x<<": "<<x<< endl; }
#define dbg2(x,y) { cerr<<#x<<": "<<x<<" , "<<#y<<": "<<y<<endl; }
#define dbg3(x,y,z) { cerr<<#x<<": "<<x<<" , "<<#y<<": "<<y<<" , "<<#z<<": "<<z<<endl; }
#define IOS ios::sync_with_stdio(false); cin.tie(0); cout.tie(0)
#ifndef LOCAL
#define endl '\n'
#endif
const int inf = INT_MAX;
const double eps = 0.0000001;
const double PI = acos(-1.0);
const int MOD = 1e9+7;
const int N  = 200005;

int m[1005];
int b[25][25];

signed main(){
    #ifdef LOCAL
        freopen("input.txt", "r", stdin);
        //freopen("output.txt", "w", stdout);
    #endif  
    IOS;

    int n;
    cin >> n;
    int a[n*n];
    for (int i = 0; i < n*n; ++i)
    {
    	cin >> a[i];
    	m[a[i]]++;
    }
    if(n%2==0) {
    	for (int i = 1; i < 1001; ++i)
	    {
	    	if(m[i]==0) continue;
	    	if(m[i]%4!=0) return cout << "NO",0;
	    }
    	int j = 0, k = 0;

    	for (int i = 1; i<=1000; i++) {
		    while (m[i]!=0) {
    			m[i]-=4;
    			b[j][k] = i;
    			k++;

    			if(k>=n/2) {
    				j++; k = 0;
    			}
    		}
    	}
    	for (int i = 0; i <= n/2; ++i)
    	{
    		for (int j = 0; j <= n/2; j++) {
    			b[n-i-1][n-j-1]=b[i][n-j-1]=b[n-i-1][j]=b[i][j];
    		}
    	}
    	cout << "YES\n";
    	for (int i = 0; i < n; ++i)
    	{
    		rep(j,0,n) {
    			cout << b[i][j] << " ";
    		}
    		cout << endl;
    	}
    	return 0;
    	
    }
    int ck1 = 0, ck2=0, ckk=0;
    int flag = 0;
    for (int i = 1; i < 1001; ++i)
    {
    	if(ck1>n-1 || ck2>1) return cout << "NO",0;
    	if(m[i]==0) continue;
    	if(m[i]%4==2) {
    		g:
    		if(flag==0) {
	    		b[ckk][n/2] = i;
	    		ckk++;
    		}
    		else {
    			b[n/2][ckk]=i; ckk++;
    		}
    		ck1++;
    		if(ck1>n-1 || ck2>1) return cout << "NO",0;
    		if(ckk>=n/2) {flag=1; ckk=0;}
    		m[i]-=2;
    	}
    	else if(m[i]%4==1 || m[i]%4==3) {
    		ck2++;
    		b[n/2][n/2] = i;
    		m[i]--;
    		if(m[i]%4==2) goto g;
    	}
    	if(ck1>n-1 || ck2>1) return cout << "NO",0;
    }
    int j = 0, k = 0;
	for (int i = 1; i<=1000; i++) {
	    while (m[i]!=0) {
	    	if(k==n/2 || j==n/2) goto a;
			m[i]-=4;
			b[j][k] = i;
			a:
			k++;
			if(k>n/2) {
				j++; k = 0;
			}
			if(j>n/2) goto b;
		}
	}

	b:;
	ckk=0;
	flag = 0;
	for (int i = 1; i < 1001; ++i)
	{
		if(m[i]==0) continue;
		// dbg2(i,m[i]);
	}
	for (int i = 1; i<=1000; i++) {
		while (m[i]!=0) {
    		if(flag==0) {
    			if(b[ckk][n/2]==0)
    			{
	    			b[ckk][n/2] = i;
	    			// dbg2(i,m[i]);
	    			m[i]-=2;
    			}
	    		ckk++;
    		}
    		else {
    			if(b[n/2][ckk]==0)
    			{
    				b[n/2][ckk]=i;
    				m[i]-=2;
    			}
    			ckk++;
    		}
    		// dbg2(flag,ckk);

    		if(ckk>=n/2) {flag=1; ckk=0;}
		}
	}
	for (int i = 0; i <= n/2; ++i)
	{
		for (int j = 0; j <= n/2; j++) {
			b[n-i-1][n-j-1]=b[i][n-j-1]=b[n-i-1][j]=b[i][j];
		}
	}
	cout << "YES\n";
	for (int i = 0; i < n; ++i)
	{
		rep(j,0,n) {
			cout << b[i][j] << " ";
		}
		cout << endl;
	}

    return 0;
}