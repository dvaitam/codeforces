#include <iostream>
#include <cstdio>
#include <cstring>
#include <algorithm>
#include <cmath>
#include <string>
#include <map>
#include <set>
#include <vector>
#include <queue>
#include <stack>
using namespace std;
typedef long long LL;
const int Maxn=1000000;
const LL Mod=10000007; 
const LL inf=0x3f3f3f3f3f3f3f;
struct node{
	int x,v;
};
node a1[100100]={},a2[100100]={};
int s1[1000100]={},s2[1000100]={};
int v1[1000100]={},v2[1000100]={};
bool cmp(node x,node y){
	if (x.x == y.x) return x.v > y.v;
	return x.x < y.x;
}
int main()
{
	int n, s;
	while (cin >> n >> s){
		int n1=0,n2=0;
		for (int i=1;i<=n;i++){
			int x,v,t;
			scanf("%d%d%d",&x,&v,&t); 
			if (t==1){
				n1++;
				a1[n1].x = x; a1[n1].v = v;
			}
			else{
				n2++;
				a2[n2].x = x; a2[n2].v = v;
			}
		}
		sort (a1+1, a1+n1+1, cmp);
		sort (a2+1, a2+n2+1, cmp);
		int pos = 0,sm = -1;
		double vm = 10000000.0;
		double t1 = 10000000.0;
		for (int i=1;i<=n1;i++){
			while (pos < a1[i].x){
				s1[pos] = sm;
				pos++;
			}
			if (vm > 1.0 / (a1[i].v + s) * a1[i].x){
				vm = 1.0 / (a1[i].v + s) * a1[i].x;
				sm = i;
			}
			t1 = min(t1, 1.0 / a1[i].v * a1[i].x);
		}
		while (pos <= 1000000){
			s1[pos] = sm;
			pos++;
		}
		pos = 1000000,vm = 10000000.0,sm = -1;
		double t2 = 10000000.0;
		for (int i=n2;i>=1;i--){
			while (pos > a2[i].x){
				s2[pos] = sm;
				pos--;
			}
			if (vm > 1.0 / (a2[i].v + s) * a2[i].x){
				vm = 1.0 / (a2[i].v + s) * a2[i].x;
				sm = i;
			}
			t2 = min(t2, 1.0 / a2[i].v * (Maxn - a2[i].x));
		}
		while (pos >=0){
			s2[pos] = sm;
			pos--;
		}
		double ans=10000000.0;
		for (int i=0;i<=1000000;i++){
			double tt1 = t1,tt2 = t2;
			if (s1[i]>0) tt1 = min(t1, 1.0 / (s + a1[s1[i]].v) * 1.0 / (s - a1[s1[i]].v) * i * s - 1.0 / (s + a1[s1[i]].v) * 1.0 / (s - a1[s1[i]].v) * a1[s1[i]].v * a1[s1[i]].x);
			if (s2[i]>0) tt2 = min(t2, 1.0 / (s + a2[s2[i]].v) * 1.0 / (s - a2[s2[i]].v) * (Maxn - i) * s - 1.0 / (s + a2[s2[i]].v) * 1.0 / (s - a2[s2[i]].v) * a2[s2[i]].v * (Maxn - a2[s2[i]].x));
			ans = min(ans, max(tt1, tt2));
		}
		printf ("%.12f\n",ans);
	}
	return 0; 
}
/*

*/