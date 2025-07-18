/*************************************************************************
    > File Name: A2.cpp
    > Author: PumpkinYing
    > Created Time: 2019/7/5 23:50:15
 ************************************************************************/

#include <iostream>
#include <stdio.h>
#include <algorithm>
#include <string>
#include <string.h>
#include <vector>
#include <queue>
#include <map>
#include <stack>
#include <cmath>
using namespace std;

#define mem(a,b) memset(a,b,sizeof(a))
typedef long long ll;

const int maxn = 1010;
vector<int> to[maxn];
int de[maxn];

void addEdge(int u,int v) {
	to[u].push_back(v);
	to[v].push_back(u);
	de[u]++;
	de[v]++;
}

int findLev(int u,int pre) {
	if(de[u] == 1) return u;
	for(auto v:to[u]) {
		if(v == pre) continue;
		return findLev(v,u);
	}
}

vector<int> work(int u,int v) {
	vector<int> ret;
	if(de[u] == 1) {
		ret.push_back(u);
		return ret;
	}
	int first = -1,second = -1;
	for(auto x:to[u]) {
		if(x == v) continue;
		else {
			if(first == -1) first = findLev(x,u);
			else if(second == -1) second = findLev(x,u);
			else break;
		}
	}
	ret.push_back(first);
	ret.push_back(second);
	return ret;
}

struct Edge {
	int u,v,val;
	Edge() {}
	Edge(int _u,int _v,int _val) : u(_u),v(_v),val(_val) {}
}es[maxn];

struct Op {
	int u,v,val;
	Op() {}
	Op(int _u,int _v,int _val) : u(_u),v(_v),val(_val) {}
}ops[maxn*10];

int tot = 0,cnt = 0;
int main() {
	int n;
	scanf("%d",&n);
	for(int i = 0;i < n-1;i++) {
		int u,v,val;
		scanf("%d%d%d",&u,&v,&val);
		addEdge(u,v);
		es[cnt++] = Edge(u,v,val);
	}

	for(int i = 1;i <= n;i++) {
		if(de[i] == 2) {
			puts("NO");
			return 0;
		}
	}

	puts("YES");
	for(int i = 0;i < cnt;i++) {
		int u = es[i].u;
		int v = es[i].v;
		int val = es[i].val;
		vector<int> ndsU = work(u,v);
		vector<int> ndsV = work(v,u);
		if(ndsU.size() == 2) {
			ops[tot++] = Op(ndsU[0],ndsU[1],-val/2);
			if(ndsV.size() == 2) {
				ops[tot++] = Op(ndsV[0],ndsV[1],-val/2);
				ops[tot++] = Op(ndsU[0],ndsV[1],val/2);
				ops[tot++] = Op(ndsU[1],ndsV[0],val/2);
			}
			else {
				ops[tot++] = Op(ndsU[0],ndsV[0],val/2);
				ops[tot++] = Op(ndsU[1],ndsV[0],val/2);
			}
		}
		else {
			if(ndsV.size() == 2) {
				ops[tot++] = Op(ndsV[0],ndsV[1],-val/2);
				ops[tot++] = Op(ndsU[0],ndsV[1],val/2);
				ops[tot++] = Op(ndsU[0],ndsV[0],val/2);
			}
			else {
				ops[tot++] = Op(ndsU[0],ndsV[0],val);
			}

		}

	}
	printf("%d\n",tot);
	for(int i = 0;i < tot;i++) printf("%d %d %d\n",ops[i].u,ops[i].v,ops[i].val);

    return 0;
}