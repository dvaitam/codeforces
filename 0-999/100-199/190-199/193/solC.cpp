//Written by Zhu Zeqi

//Come on,baby

//Hack,please

#include<cmath>

#include<math.h>

#include<ctype.h>

#include<algorithm>

#include<bitset>

#include<cassert>

#include<cctype>

#include<cerrno>

#include<cfloat>

#include<ciso646>

#include<climits>

#include<clocale>

#include<complex>

#include<csetjmp>

#include<csignal>

#include<cstdarg>

#include<cstddef>

#include<cstdio>

#include<cstdlib>

#include<cstring>

#include<ctime>

#include<cwchar>

#include<cwctype>

#include<deque>

#include<exception>

#include<fstream>

#include<functional>

#include<iomanip>

#include<ios>

#include<iosfwd>

#include<iostream>

#include<istream>

#include<iterator>

#include<limits>

#include<list>

#include<locale>

#include<map>

#include<memory>

#include<new>

#include<numeric>

#include<ostream>

#include<queue>

#include<set>

#include<sstream>

#include<stack>

#include<stdexcept>

#include<streambuf>

#include<string>

#include<typeinfo>

#include<utility>

#include<valarray>

#include<vector>

#include<string.h>

#include<stdlib.h>

#include<stdio.h>

#define ll   long long

#define pb push_back

#define mp make_pair

#define F first

#define S second

#define pii pair<int,int>

#define vi vector<int>

#define MAX 100000000000000000

#define MOD 1000000007

#define PI 3.141592653589793238462

#define INF 1000000000

using namespace std;

int zh[4][4],zhu[4],t;

char str[4][400005];

void los(){

	cout<<-1<<endl;

	exit(0);

}

void sing(int i,int m){

	for(int j=0;j<4;j++){

		if(j==i)

		for(int k=0;k<m;k++)

		str[j][t+k]='b';

		else{

			for(int k=0;k<m;k++)

			str[j][t+k]='a';

			zh[i][j]-=m;

			zh[j][i]-=m;

		}

	}

	t+=m;

}

void solve(int i,int m){

	if(m<0)

	los();

	for(int j=0;j<4;j++){

		if(j==0 || j==i)

		for(int k=0;k<m;k++)

		str[j][t+k]='b';

		else

		for(int k=0;k<m;k++)

		str[j][t+k]='a';

	}

	t+=m;

}

int main(){

	//freopen("input.in","r",stdin);

	//freopen("output.out","w",stdout);

	for(int i=0;i<=2;i++)

	for(int j=i+1;j<=3;j++){

		cin>>zh[i][j];

		zh[j][i]=zh[i][j];

		zhu[i]+=zh[i][j];

		zhu[j]+=zh[j][i];

	}

	int res=zhu[0];

	for(int i=0;i<4;i++)

	res=min(res,zhu[i]);

	for(int i=0;i<4;i++){

		if((zhu[i] & 1)!=(res & 1))

		los();

		sing(i,(zhu[i]-res)/2);

	}

	int as=zh[0][1]+zh[0][2]+zh[0][3];

	if(as & 1)

	los();

	solve(3,(zh[0][1]+zh[0][2]-zh[0][3])/2);

	solve(2,(zh[0][1]+zh[0][3]-zh[0][2])/2);

	solve(1,(zh[0][3]+zh[0][2]-zh[0][1])/2);

	cout<<t<<endl;

	for(int i=0;i<4;i++){

		str[i][t]=0;

		cout<<str[i]<<endl;

	}

	return 0;

}