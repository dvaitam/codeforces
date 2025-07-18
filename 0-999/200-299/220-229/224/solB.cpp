#include <iostream>
#include <sstream>
#include <cstdio>
#include <cstring>
#include <cmath>
#include <vector>
#include <cassert>
#include <list>
#include <queue>
#include <algorithm>
#include <set>
#include <map>
#define REP(i,n) for(int i=0; i<(int)n; i++)
#define FOREACH(it,v) for(typeof((v).begin()) it=(v).begin(); it!=(v).end(); ++it)
using namespace std;
const int INF=1000000000;


string str(const string& s) {
	return s;
}

string str(int x) {
	ostringstream ss;
	ss<<x;
	return ss.str();
}

template<typename T>
string str(T v) {
	ostringstream ss;
	ss<<"[";
	FOREACH(x,v) ss<<str(*x)<<",";
	ss<<"]";
	return ss.str();
}



int main()
{
	ios_base::sync_with_stdio(false);
	int N,K;
	cin>>N>>K;
	vector<int> data(N);
	vector<int> veces(100001);
	int difs=0;
	REP(i,N) {
		cin>>data[i];
		if(++veces[data[i]]==1)
			difs++;
	}
	
	if(difs<K) {
		cout<<"-1 -1"<<endl;
		return 0;
	}
	
	int l, r;
	for(l=0; l<N; l++)
		if(--veces[data[l]]==0) 
			if(--difs<K) {
				difs++;
				veces[data[l]]++;
				break;
			}
		
	
	
	for(r=N-1; r>=0; r--) 
		if(--veces[data[r]]==0)
			if(--difs<K) {
				difs++;
				veces[data[r]]++;
				break;
			}
	
	//veces.resize(10);
	//cout<<str(veces)<<endl;
	cout<<l+1<<" "<<r+1<<endl;
}