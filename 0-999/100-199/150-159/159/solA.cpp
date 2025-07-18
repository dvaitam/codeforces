/* Author haleyk10198 */

/* �@��:  haleyk10198 */

#include <iostream>

#include <fstream>

#include <sstream>

#include <cstdlib>

#include <cstdio>

#include <vector>

#include <map>

#include <queue>

#include <cmath>

#include <algorithm>

#include <cstring>

#include <iomanip>

#include <ctime>

#include <string>

#include <set>

#include <stack>



#define MOD 1000000007

#define INF 2147483647

#define PI 3.1415926535897932384626433

#define ll long long

#define pii pair<int,int>

#define mp(x,y) make_pair((x),(y))



using namespace std;



vector<string> nlist;

map<string,int> sconv;

map<pii,int> pconv;

set<pii> out;

vector<set<int> > lg;



int main(){

	ios_base::sync_with_stdio(false);

	int n,d;

	cin>>n>>d;

	cin.ignore();

	for(int i=0;i<n;i++){

		string str;

		getline(cin,str,'\n');

		string n1,n2;

		int t,pos1,pos2;

		pos1=str.find(' ');

		pos2=str.find(' ',pos1+1);

		n1=str.substr(0,pos1);

		n2=str.substr(pos1+1,pos2-pos1-1);

		t=stoi(str.substr(pos2+1));

		if(sconv.count(n1)==0){

			sconv[n1]=nlist.size();

			nlist.push_back(n1);

		}

		if(sconv.count(n2)==0){

			sconv[n2]=nlist.size();

			nlist.push_back(n2);

		}

		int p1=sconv[n1],p2=sconv[n2];

		if(pconv.count(mp(p1,p2))==0){

			pconv[mp(p1,p2)]=lg.size();

			lg.push_back(set<int>());

		}

		if(pconv.count(mp(p2,p1))==0){

			pconv[mp(p2,p1)]=lg.size();	

			lg.push_back(set<int>());

		}

		lg[pconv[mp(p1,p2)]].insert(t);

		set<int>::iterator it=lg[pconv[mp(p2,p1)]].lower_bound(t-d);

		if(it!=lg[pconv[mp(p2,p1)]].end()){

			if(*it!=t){

				if(p1>p2)

					swap(p1,p2);

				out.insert(mp(p1,p2));

			}

		}		

	}

	cout<<out.size()<<endl;

	for(auto x:out)

		cout<<nlist[x.first]<<" "<<nlist[x.second]<<endl;

	return 0;

}