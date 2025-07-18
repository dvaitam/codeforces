#define Federico using
#define Javier namespace
#define Pousa std;
#include <iostream>
#include <fstream>
#include <string>
#include <set>
#include <vector>
#include <map>
#include <algorithm>
#include <cstdio>
#include <cstdlib>
#include <cmath>
#include <stack>
#include <queue>
#include <cstring>
#include <sstream>


Federico Javier Pousa

int in(){int r=0,c;for(c=getchar();c<=32;c=getchar());if(c=='-') return -in();for(;c>32;r=(r<<1)+(r<<3)+c-'0',c=getchar());return r;}

#define INF 100000000000000000LL

typedef pair<int,int> pii;
typedef pair<long long, long long> lii;
long long int cFijo = 1000000;
vector<pii> ejes[100005];
vector<long long int> dist;
vector<int> padres;
int N, M, X, Y, Z, rotos;

void dij(int s){
	set<lii> Q;
	dist[s] = 0;
	Q.insert(lii(0,s));
	
	while(!Q.empty()){
		lii top = *Q.begin();
		Q.erase(Q.begin());
		int v = top.second;
		
		for(int i=0; i<(int)ejes[v].size(); ++i){
			int v2 = ejes[v][i].first;
			long long cost = cFijo + (ejes[v][i].second?0:1);
			
			if(dist[v2] > dist[v]+cost){
				if(dist[v2] != INF){
					Q.erase(Q.find(lii(dist[v2], v2)));
				}
				dist[v2] = dist[v] + cost;
				Q.insert(lii(dist[v2], v2));
				padres[v2] = v;
			}
		}
	}
	
	return;
}


int main(){
	N = in();
	M = in();
	dist = vector<long long int>(N+5, INF); 
	padres = vector<int>(N+5, -1);
	for(int i=0; i<M; ++i){
		X = in()-1;
		Y = in()-1;
		Z = in();
		ejes[X].push_back(pii(Y,Z));
		ejes[Y].push_back(pii(X,Z));
		if(!Z)rotos++;
	}
	
	dij(0);
	
	set<pii> enCamino;
	int act = N-1;
	int padre = padres[act];
	while(padre!=-1){
		enCamino.insert(pii(min(act,padre),max(act,padre)));
		act = padre;
		padre = padres[padre];
	}
	
	//~ cout << M-rotos-(dist[N-1]/cFijo)+2*(dist[N-1]%cFijo) << endl;
	printf("%d\n", M-rotos-(dist[N-1]/cFijo)+2*(dist[N-1]%cFijo));
	int vec;
	for(int i=0; i<N; ++i){
		for(int j=0; j<(int)ejes[i].size(); ++j){
			vec = ejes[i][j].first;
			if(i>vec)continue;
			if(enCamino.find(pii(i,vec))!=enCamino.end()){
				if(ejes[i][j].second==0){
					printf("%d %d %d\n", i+1, vec+1, 1);
				}
			}else{
				if(ejes[i][j].second==1){
					printf("%d %d %d\n", i+1, vec+1, 0);
				}
			}
		}
	}
	return 0;
}