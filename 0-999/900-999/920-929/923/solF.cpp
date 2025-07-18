#include <map>

#include <set>

#include <queue>

#include <cmath>

#include <bitset>

#include <cstdio>

#include <cstring>

#include <iostream>

#include <algorithm>

#define pii pair <int , int>

#define mp make_pair

#define fs first

#define sc second

using namespace std;

typedef long long LL;

typedef unsigned long long ULL;



//const int Mxdt=100000; 

//static char buf[Mxdt],*p1=buf,*p2=buf;

//#define getchar() p1==p2&&(p2=(p1=buf)+fread(buf,1,Mxdt,stdin),p1==p2)?EOF:*p1++;



template <typename T>

void read(T &x) {

	T f=1;x=0;char s=getchar();

	while(s<'0'||s>'9') {if(s=='-') f=-1;s=getchar();}

	while(s>='0'&&s<='9') {x=(x<<3)+(x<<1)+(s-'0');s=getchar();}

	x *= f;

}



template <typename T>

void write(T x , char s='\n') {

	if(!x) {putchar('0');putchar(s);return;}

	if(x<0) {putchar('-');x=-x;}

	T tmp[25]={},t=0;

	while(x) tmp[t++]=x%10,x/=10;

	while(t-->0) putchar(tmp[t]+'0');

	putchar(s); 

}



const int MAXN = 1e5 + 5;



int head[MAXN] , to[MAXN << 1] , nxt[MAXN << 1] , cnt , deg[MAXN] , n;

void add(int u , int v) {nxt[++cnt] = head[u];head[u] = cnt;to[cnt] = v;} 



set <int> S[MAXN];

set <pii> D1 , D2;

set <int> res1 , res2;

queue <int> q1 , q2;



void pre(int x , int las) {

	for (int i = head[x]; i; i = nxt[i]) {

		int v = to[i];

		if(v == las) continue;

		pre(v , x);

		if(deg[v] == 1) S[x].insert(v);

	} 

	if(S[x].size()) {

		if(x <= n) q1.push(x);

		else q2.push(x);

	}

}



int p[MAXN] , v1[MAXN] , v2[MAXN] , t , vis[MAXN];

pii kkkk;



bool check() {

	for (int i = 1; i <= t; ++i) {

		int x = v1[i];

		for (int k = head[x]; k; k = nxt[k]) {

			int v = to[k];

			vis[v] ++;

		}

		

		x = p[x];

		for (int k = head[x]; k; k = nxt[k]) {

			int v = to[k];

			vis[p[v]] ++;

		}

		

		int fl = 1;

		for (int j = 1; j <= t; ++j) if(vis[v1[j]] > 1) fl = 0 , kkkk = mp(v1[i] , v1[j]);

		

		x = v1[i];

		for (int k = head[x]; k; k = nxt[k]) {

			int v = to[k];

			vis[v] --;

		}

		

		x = p[x];

		for (int k = head[x]; k; k = nxt[k]) {

			int v = to[k];

			vis[p[v]] --;

		}

		if(!fl) return false;

	}

	return true;

} 



bool dfs(int x) {

	if(x == t + 1) {

		if(check()) return true;

		return false;

	}

	

	for (int i = 1; i <= t; ++i) if(!p[v2[i]]) {

		p[v1[x]] = v2[i];

		p[v2[i]] = v1[x];

		if(dfs(x + 1)) return true;

		p[v1[x]] = 0;

		p[v2[i]] = 0; 

	}

	return false;

}



void solve(int n) {

	if(n <= 5) {

		t = n;

		for (int i = 1; i <= n; ++i) v1[i] = *res1.begin() , res1.erase(v1[i]);

		for (int i = 1; i <= n; ++i) v2[i] = *res2.begin() , res2.erase(v2[i]);

		dfs(1);

		return;

	}

	if((*D2.rbegin()).fs == n - 2) D1.swap(D2) , q1.swap(q2) , res1.swap(res2);

	if((*D1.rbegin()).fs == n - 2) { // one-star

		int u = q1.front();q1.pop();

		int w = q1.front();q1.pop();

		if(S[u].size() < S[w].size()) swap(u , w);

		int v = *S[w].begin();

		

		res1.erase(u),res1.erase(w),res1.erase(v);

		

		int x = q2.front();

		p[v] = x , p[x] = v;

		res2.erase(x);

		x = *S[x].begin();q2.pop();

		p[u] = x , p[x] = u;

		res2.erase(x);

		x = q2.front();q2.pop();

		x = *S[x].begin();

		p[w] = x , p[x] = w;

		res2.erase(x);

		

		while(res1.size()) {

			int u = *res1.begin() , v = *res2.begin();

			p[u] = v , p[v] = u;

			res1.erase(u) , res2.erase(v);

		}

		return;

	}

	

	int u[4] , fu[4];

	

	for (int k = 0; k < 2; ++k) {

		fu[k] = q1.front();q1.pop();

		D1.erase(mp(deg[fu[k]] , fu[k]));

		u[k] = *S[fu[k]].begin();if(k != 0 && u[k] == fu[k - 1]) u[k] = *S[fu[k]].rbegin();

		S[fu[k]].erase(u[k]);

		deg[fu[k]] --;res1.erase(u[k]);

		if(deg[fu[k]] > 1) {

			D1.insert(mp(deg[fu[k]] , fu[k]));

			if(S[fu[k]].size()) q1.push(fu[k]);

		}

		else {

			

			int vv = 0;

			for (int i = head[fu[k]]; i; i = nxt[i]) {

				int v = to[i];

				if(p[v] || v == u[k]) continue;

				else {

					vv = v;

					S[v].insert(fu[k]);

					if(S[v].size() == 1u) q1.push(v); 

				}

			}

			

			if(q1.size() == 1) { // two-star

				res1.insert(u[k]);deg[fu[k]] ++;

				S[fu[k]].insert(u[k]);

				if(vv) S[vv].erase(fu[k]);

				D1.insert(mp(deg[fu[k]] , fu[k]));

				q1.push(fu[k]);

				k --;

				continue;

			}

		}

	}

	

	for (int k = 2; k < 4; ++k) {

		fu[k] = q2.front();q2.pop();

		D2.erase(mp(deg[fu[k]] , fu[k]));

		u[k] = *S[fu[k]].begin();if(k != 0 && u[k] == fu[k - 1]) u[k] = *S[fu[k]].rbegin();

		S[fu[k]].erase(u[k]);

		deg[fu[k]] --;res2.erase(u[k]);

		if(deg[fu[k]] > 1) {

			D2.insert(mp(deg[fu[k]] , fu[k]));

			if(S[fu[k]].size()) q2.push(fu[k]);

		}

		else {

			int vv = 0;

			for (int i = head[fu[k]]; i; i = nxt[i]) {

				int v = to[i];

				if(p[v] || v == u[k]) continue;

				else {

					vv = v;

					S[v].insert(fu[k]);

					if(S[v].size() == 1u) q2.push(v); 

				}

			}

			

			if(q2.size() == 1) { // two-star

				res2.insert(u[k]);deg[fu[k]] ++;

				S[fu[k]].insert(u[k]);

				if(vv) S[vv].erase(fu[k]);

				D2.insert(mp(deg[fu[k]] , fu[k]));

				q2.push(fu[k]);

				k --;

				continue;

			}

		}

	}

	

	p[u[0]] = u[2] , p[u[2]] = u[0];

	p[u[1]] = u[3] , p[u[3]] = u[1]; 

	solve(n - 2);

	if(p[fu[0]] == fu[2] || p[fu[1]] == fu[3]) swap(p[u[0]] , p[u[1]]) , swap(p[u[2]] , p[u[3]]); 

	return;

	

}



int main() {

//	freopen("1.in" , "r" , stdin);

//	freopen("1.out" , "w" , stdout);

	read(n);

	for (int i = 1; i < n; ++i) {

		int u , v;

		read(u),read(v);

		add(u , v) , add(v , u); 

		deg[u] ++ , deg[v] ++;

	}

	for (int i = 1; i < n; ++i) {

		int u , v;

		read(u),read(v);

		deg[u] ++ , deg[v] ++;

		add(u , v) , add(v , u); 

	}

	

	for (int i = 1; i <= n; ++i) if(deg[i] != 1) {

		pre(i , 0);

		break;

	}

	for (int i = n + 1; i <= 2 * n; ++i) if(deg[i] != 1) {

		pre(i , 0);

		break;

	}

	

	for (int i = 1; i <= n; ++i) res1.insert(i) , res2.insert(i + n) , D1.insert(mp(deg[i] , i)) , D2.insert(mp(deg[i + n] , i + n));

	

	if(q1.size() <= 1 || q2.size() <= 1) {

		puts("No");

		return 0;

	}

	

	puts("Yes");

	solve(n);

	

	for (int i = 1; i <= n; ++i) write(p[i] , ' ') , v1[i] = i , v2[i] = n + i;

	t = n;

//	freopen("2.out" , "w" , stdout);

//	cerr << endl;

//	cerr << check() << endl;

	

	return 0;

}