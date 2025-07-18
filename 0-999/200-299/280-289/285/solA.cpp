#include <valarray>
#include <complex>
#include <string>
#include <functional>
#include <iterator>
#include <algorithm>
#include <vector>
#include <stack>
#include <queue>
#include <set>
#include <map>
#include <list>
#include <bitset>
#include <locale>
#include <memory>
#include <stdexcept>
#include <exception>
#include <sstream>
#include <iomanip>
#include <fstream>
#include <iostream>
#include <cstdarg>
#include <cstdlib>
#include <ctime>
#include <cmath>
#include <cctype>
#include <cstring>
#include <cstdio>
using namespace std;
int N,K;
void islem(){
	int i,j;
	for( i=N ; i>N-K ; i-- )	cout<<i<<' ';
	for( j=1 ; j<=i ; j++ )	cout<<j<<' ';
	cout<<endl;
}

int main(){
	cin>>N>>K;
	 islem();
	return 0;
}