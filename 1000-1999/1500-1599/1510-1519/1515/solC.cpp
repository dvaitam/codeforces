#include <algorithm>

#include <bits/stdc++.h>

#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp>

#include <cstdint>

#include <iostream>

#include <climits>

#include <locale>

#include <memory>

#include <string>

#include <utility>

using namespace std;

using namespace __gnu_pbds;

#define int long long

#define mod 1000000007

#define inf 10000000000000001

#define pb push_back

#define endl '\n'

#define vi vector<int>

#define all(x) (x).begin(),(x).end()

#define rall(x) (x).rbegin(),(x).rend()

#define sz(x) (int)(x).size()

#define vvi vector<vector<int>>

#define pii pair<int,int>





template<class T> 

using ordered_set = tree<T, null_type,less<T>, rb_tree_tag, tree_order_statistics_node_update>;



template<class T> 

using ordered_multiset = tree<T, null_type, less_equal<T>, rb_tree_tag, tree_order_statistics_node_update>;



int pcnt(int x) {

  int ans=0;

  while(x) {

    ans+=(x&1);

    x>>=1;

  }

  return ans;

}



int binpow(int a,int b) {

  a%=mod;

  int res=1;

  while(b>0) {

    if(b&1)

      res=res*a%mod;

    a=a*a%mod;

    b>>=1;

  }

  return res;

}



int nCr(int n,int k) {

  if(k>n) return 0;

  int numerator=1;

  for(int i=0; i<k; i++) {

    numerator=(numerator*(n-i))%mod;

  }

  int denominator=1;

  for (int i=1;i<=k;i++) {

    denominator=(denominator*i)%mod;

  }

  return (numerator*binpow(denominator,mod-2))%mod;

}

 

int nPr(int n,int k) {

  int ans=1;

  for(int i=n-k+1;i<=n;i++) {

    ans*=i;

    ans%=mod;

  }

  return ans;

}

 

string longdif(string str1, string str2) {

  string str="";

  int n1=str1.length(), n2=str2.length();

  int diff=n1-n2;

  int carry=0;

  for(int i=n2-1;i>=0;i--) {

    int sub=((str1[i+diff]-'0')-(str2[i] - '0')-carry);

    if(sub < 0) {

      sub=sub+10;

      carry=1;

    }

    else carry=0;

    str+=sub+'0';

  }

  for(int i=n1-n2-1;i>=0;i--) {

    if(str1[i]=='0' && carry) {

      str+='9';

      continue;

    }

    int sub=((str1[i]-'0')-carry);

    if(i>0||sub>0)

      str+=sub+'0';

    carry=0;

  }

  reverse(str.begin(), str.end());

  return str;

}

 

string findSum(string str1, string str2){

  if (str1.length() > str2.length()) swap(str1, str2);

  string str = "";

  int n1 = str1.length(), n2 = str2.length();

  reverse(str1.begin(), str1.end());

  reverse(str2.begin(), str2.end());

  int carry = 0;

  for (int i=0; i<n1; i++){

      int sum = ((str1[i]-'0')+(str2[i]-'0')+carry);

      str.push_back(sum%10 + '0');

      carry=sum/10;

  }

  for (int i=n1; i<n2; i++){

    int sum=((str2[i]-'0')+carry);

    str.push_back(sum%10+'0');

    carry = sum/10;

  }

  if(carry) str.push_back(carry+'0');

  reverse(str.begin(), str.end());

  return str;

}

 

void dfs(vector<vi> &adj,vi &vis,int u,int c) {

  vis[u]=c;

  for(auto v:adj[u]) {

    if(vis[v]!=c) {

      dfs(adj,vis,v,c);

    }

  }

}



template<class T> struct iseg {

	const T ID=0; 

  T comb(T a, T b) { 

    return a+b;

  }

	int n;

  vector<T> seg;

	void init(int _n) {

    n=_n;

    seg.assign(2*n,ID);

  }

  void finit() {

    for(int i=n-1;i>=1;i--) 

      seg[i]=comb(seg[i<<1],seg[(i<<1)+1]);

  }

	void pull(int p) {

    seg[p]=comb(seg[2*p],seg[2*p+1]); 

  }

	void upd(int p,T val) {

		seg[p+=n]=val; 

    for(p/=2;p;p/=2) 

      pull(p);

  }

	T query(int l,int r) {	

		T ra=ID,rb=ID;

		for(l+=n,r+=n+1;l<r;l/=2,r/=2) {

			if(l&1) ra=comb(ra,seg[l++]);

			if(r&1) rb=comb(seg[--r],rb);

		}

		return comb(ra,rb);

	}

};



void solve() {

  int n,m,x;

  cin>>n>>m>>x;

  vector<pair<int,int>> a(n);

  for(int i=0;i<n;i++) {

    int x;

    cin>>x;

    a[i]={x,i};

  }

  sort(all(a));

  vi b(n,0);

  cout<<"YES"<<endl;

  for(int i=0;i<n;i++) {

    b[a[i].second]=i%m+1;

    //cout<<i%m+1<<' ';

  }

  for(auto u:b) {

    cout<<u<<' ';

  }

  cout<<endl;

}

 

signed main() {

  ios_base::sync_with_stdio(false);

  cin.tie(NULL);

  cout.tie(NULL);

  int tt=1;

  cin>>tt;

  while(tt--) {

    solve();

  }

}