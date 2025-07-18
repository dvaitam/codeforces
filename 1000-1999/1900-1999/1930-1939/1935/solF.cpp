#include<stdio.h>
#include<string.h>
#include<stdlib.h>

typedef long long LL;
#include <map>
#include <vector>
#include <queue>
#include <deque>
#include <set>
#include <stack>
#include <algorithm>
#include <array>
#include <unordered_set>
#include <unordered_map>
#include <string>
using namespace std;

#define P1 972663749
#define P2 911382323
#define MOD 998244353

bool rcmp(int a, int b) { return a>b; }
struct VNode {
	int v;
	map<int, int> *p;
	bool operator<(const VNode& b) const {
		return v>b.v;
	}
};

vector<int> nei[1<<19];
map<int, int> xx[1<<19];
map<int, int> *pxx[1<<19];
VNode vs[1<<19];
map<int, int> *mergeit(map<int, int>* a, map<int, int>* b) {
	if (a==NULL) return b;
	if (b==NULL) return a;
	if (a->size()<b->size()) swap(a, b);
	// printf("merge"); for (auto x=a->begin(); x!=a->end(); x++) printf("[%d,%d]", x->first, x->second);
	// printf("with"); for (auto x=b->begin(); x!=b->end(); x++) printf("[%d,%d]", x->first, x->second);
	int s, e;
	for (auto x = b->begin(); x!=b->end(); x++) {
		s=x->first; e=x->second;
		auto y = a->lower_bound(s);
		if (y==a->begin()) {
			if (e+1==y->first) {
				e=y->second;
				a->erase(y);
			} (*a)[s]=e;
		} else {
			auto z=y; z--;
			if (z->second==s-1) {
				s=z->first;
				if (y!=a->end()&&e+1==y->first) {
					e=y->second;
					a->erase(y);
				} (*a)[s]=e;
			} else {
				if (y!=a->end()&&e+1==y->first) {
					e=y->second;
					a->erase(y);
				} (*a)[s]=e;
			}
		}
	}
	// printf("==>"); for (auto x=a->begin(); x!=a->end(); x++) printf("[%d,%d]", x->first, x->second); printf("\n");
	return a;
}
vector<pair<int,int>> rs[1<<19];
int gn;
int ss[1<<19], es[1<<19];
int gl;
void dfs(int a, int p) {
	ss[a]=gl++;
	for (auto b: nei[a]) {
		if (b==p) continue;
		dfs(b, a);
	}
	es[a]=gl-1;
}

void build(int r, int rp) {
	int k=0, i, j, i1, i2, v;
	xx[r].clear();
	rs[r].clear();
	xx[r][r]=r;
	map<int, int> *p1=NULL, *p=NULL, *p2=NULL;
	for (auto b: nei[r]) {
		if (b==rp) continue;
		build(b, r);
	}
	for (auto b: nei[r]) {
		if (b==rp) continue;
	       	p=pxx[b];
	       	vs[k].v=p->rbegin()->second;
		vs[k].p=p;
	       	k++;
	} sort(vs, vs+k);
	p = NULL;
	for (i=0; i<k; i++) {
		v = vs[i].v;
		j = v+1;
		if (j<=gn&&(ss[j]<ss[r]||ss[j]>es[r])) {
			// connect to outer
			p2 = mergeit(p2, vs[i].p);
			rs[r].push_back({v, v+1});
		} else {
			if (p2) {
				auto z = p2->find(j);
				if (z!=p2->end()) {
					p2=mergeit(p2, vs[i].p);
					rs[r].push_back({v, v+1});
					continue;
				}
			}
			if (p==NULL) p=vs[i].p;
			else {
				if (j==r) {
					p1=vs[i].p;
					for (i++; i<k; i++) {
						v=vs[i].v;
					       	rs[r].push_back({v, v+1});
						j=v+1;
						if (j<=gn&&(ss[j]<ss[r]||ss[j]>es[r])) {
							p2 = mergeit(p2, vs[i].p);
						} else {
							if (p2) {
								auto z = p2->find(j);
								if (z!=p2->end()) {
									p2=mergeit(p2, vs[i].p);
									continue;
								}
							}
							auto z = p1->find(vs[i].p->rbegin()->second+1);
							if (z==p1->end()) {
								p = mergeit(p, vs[i].p);
							} else {
								p1 = mergeit(p1, vs[i].p);
							}
						}
					}
					break;
				} else {
					rs[r].push_back({v, v+1});
					p = mergeit(p, vs[i].p);
				}
			}
		}
	}
	if (p==NULL) {
		p = mergeit(p2, &xx[r]);
	} else {
		if (p1!=NULL) {
			// find an edge for p1
			// printf("%d:", r); for (auto z=p1->begin(); z!=p1->end(); z++) printf("[%d,%d]", z->first, z->second); printf("\n");
			auto x = p1->rbegin();
			if (x->first==1) {
				rs[r].push_back({x->second, x->second+2});
				i = x->second+2;
			} else {
				rs[r].push_back({x->first-1, x->first});
				i = x->first-1;
			}
			char fit=0;
			if (p2) {
				auto z = p2->upper_bound(i);
				if (z!=p2->begin()) {
					z--;
					if (z->second>=i) {
						fit=1;
					}
				}
			}
			if (fit||ss[i]<ss[r]||ss[i]>es[r]) {
				p2 = mergeit(p2, p1); p1=NULL;
			} else {
				p = mergeit(p, p1); p1=NULL;
			}
		}
		p = mergeit(p, &xx[r]);
		i=j=-1;
		for (auto x=p->begin(); x!=p->end(); x++) {
			if (x->first==x->second&&x->first==r) continue;
			i1=x->first-1; i2=x->first; if (i2==r) i2++;
			if (i1>=1&&i2<=gn) {
				if (i==-1||i2-i1==1) { i=i1; j=i2; }
			}
			i1=x->second; i2=x->second+1; if (i1==r) i1--;
			if (i1>=1&&i2<=gn) {
				if (i==-1||i2-i1==1) { i=i1; j=i2; }
			}
			if (j-i==1) break;
		}
		if (i!=-1) rs[r].push_back({i, j});
		p = mergeit(p, p2);
	}
	// printf("%d:", r); for (auto x=p->begin(); x!=p->end(); x++) printf("[%d,%d]", x->first, x->second); printf("\n");
	pxx[r]=p;
}
char cmk[1280];
vector<int> nnei[1280];
int que[1280];
char check(int n, vector<pair<int, int>>& es, int z) {
	int i, a, b;
	int h=0, t=0;
	if (nei[z].size()-1!=es.size()) return 0;
	for (i=1; i<=n; i++) cmk[i]=0;
	for (i=1; i<=n; i++) nnei[i].clear();
	for (auto x: es) {
		a=x.first; b=x.second;
		if (a==z) return 0;
		if (b==z) return 0;
		nnei[a].push_back(b);
		nnei[b].push_back(a);
	}
	a=1; if (z==a) a++;
	cmk[a]=1; que[h++]=a;
	while(t<h) {
		a=que[t++];
		for (auto b: nei[a]) {
			if (b==z) continue;
			if (cmk[b]) continue;
			cmk[b]=1; que[h++]=b;
		}
		for (auto b: nnei[a]) {
			if (b==z) continue;
			if (cmk[b]) continue;
			cmk[b]=1; que[h++]=b;
		}
	}
	if (h==n-1) return 1;
	return 0;
}
int main() {
	int n, i, a, m, c, b;
	int tc; scanf("%d", &tc); while(tc--) {
		scanf("%d", &n); gn=n;
		for (i=1; i<=n; i++) nei[i].clear();
		for (i=1; i<n; i++) {
			scanf("%d %d", &a, &b);
			nei[a].push_back(b);
			nei[b].push_back(a);
		}
		gl=0; dfs(1, -1);
		build(1, -1);
		for (a=1; a<=n; a++) {
			m=0; c=0;
			for (auto x: rs[a]) {
				m+=abs(x.first-x.second);
				c++;
			}
			/*
			if (check(n, rs[a], a)==0) {
				for (a=1; a<=n; a++) {
					printf("%d: ", a); for (auto b: nei[a]) printf(" %d", b); printf("\n");
				}
			       	printf("error\n");
			       	exit(1);
		       	}
			if (tc>10) continue;
			*/
			printf("%d %d\n", m, c);
			for (auto x: rs[a]) printf("%d %d\n", x.first, x.second);
			printf("\n");
		}
	}
	return 0;
}