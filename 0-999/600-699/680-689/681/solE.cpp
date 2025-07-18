#include<bits/stdc++.h>

#define f first

#define s second

#define pb push_back

#define mp make_pair

#define pi acos(-1.)



using namespace std;

typedef long long ll;



int const N = (1e5) + 5;

vector<ll> x, y, r;

vector<pair<double, double>> seg;

int n;

double alpha(double dx, double dy, double dist){

	double angle = acos(dx / dist);

	if(dy < 0){

		angle += acos(-1.);

	}

	return angle;

}

int main(){

	ll x0, y0, v, T;

	cin >> x0 >> y0 >> v >> T;

	cin >> n;

	for(int i = 0; i < n; ++i){

		int xi, yi, ri;

		scanf("%d %d %d", &xi, &yi, &ri);

		x.pb(xi), y.pb(yi), r.pb(ri);

	}

	n = x.size();

	double R = min(v * T, 6 * 1000000000LL);

	for(int i = 0; i < n; ++i){

		double dx = x[i] - x0, dy = y[i] - y0;

		double dist = dx * dx + dy * dy;

		if((x[i] - x0) * (x[i] - x0) + (y[i] - y0) * (y[i] - y0) <= r[i] * r[i]){

			cout << 1;

			return 0;

		}

		dist = sqrt(dist);

		if(R + r[i] <= dist || r[i] == 0){

			continue;

		}

		double ang1, ang2;

		double angl = atan2(dx, dy); //alpha(dx, dy, dist);

		//cout << setprecision(9) << atan2(dx, dy) << " " << alpha(dx, dy, dist) << endl;

		double l = dist * dist - r[i] * r[i];

		double al = asin(r[i] / dist);

		if(R * R + r[i] * r[i] < dist * dist){

			double Cos = (R * R + dist * dist - r[i] * r[i]) / (2. * R * dist);

			al = acos(Cos);

		}

		/*if((l <= (R * R))){ // || (dist <= R)){

			double al = asin(r[i] / dist);

		}else{

			double Cos = (R * R + dist * dist - (ll)r[i] * r[i]) / (2. * R * dist);

			double al = acos(Cos);

		}*/

		ang1 = angl - al;

		ang2 = angl + al;

		if(ang1 < -pi){

			seg.push_back(make_pair(-pi, ang2));

			seg.push_back(make_pair(ang1 + 2. * acos(-1.), pi));

		}else if(ang2 > pi){

			seg.push_back(make_pair(ang1, pi));

			seg.push_back(make_pair(-pi, ang2 - 2. * acos(-1.)));

		}else{

			seg.push_back(make_pair(ang1, ang2));

		}

		/*if(ang1 < 0){

			seg.push_back(make_pair(0, ang2));

			seg.push_back(make_pair(ang1 + 2. * acos(-1.), 2. * acos(-1.)));

		}else if(ang2 > 2 * acos(-1.)){

			seg.push_back(make_pair(ang1, 2. * acos(-1.)));

			seg.push_back(make_pair(0, ang2 - 2. * acos(-1.)));

		}else{

			seg.push_back(make_pair(ang1, ang2));

		}*/

	}

	sort(seg.begin(), seg.end());

	

	/*for(int i = 0; i < seg.size(); ++i){

		cout << setprecision(9) << seg[i].f << " " << seg[i].s << endl;

	}*/

	double len = 0, clen = 0;

	double en = -pi;

	for(int i = 0; i < seg.size(); ++i){

		//en = max(en, seg[i].f);

		//if(seg[i].s - en > 0) len += seg[i].s - en;

		//en = max(en, seg[i].s);

		if(en  >= seg[i].f){

			if(en < seg[i].s){

				clen += (seg[i].s - en);

				en = seg[i].s;

			}

		}else{

			len += clen;

			clen = seg[i].s - seg[i].f;

			en = seg[i].s;

		}

	}

	len += clen;

	//cout << setprecision(9) << len << endl;

	double p = len / (2 * acos(-1.));

	cout << setprecision(9) << p << endl;

	return 0;

}



/*if((l <= (R * R)) || (dist <= R)){

			double al = asin(r[i] / dist);

		}else{

			//double p = (R + r[i] + dist) / 2.;

			//double s = sqrt(p * (p - R)) * sqrt((p - r[i]) * (p - dist));

			//double al = asin(2 * s / (R * dist));

			double Cos = (R * R + dist * dist - (ll)r[i] * r[i]) / (2. * R * dist);

			double al = acos(Cos);

		}*/