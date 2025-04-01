package algorithm

import "strconv"

func muchZero(x int) string {
	var m = ""
	for i := 0; i < x; i++ {
		m += "0"
	}
	return m
}

const d = 6

func qf(a string) []int {
	var as = make([]int, 0, len(a)/8+1)
	for i := 0; i < len(a); i += 8 {
		var v = 0
		var err error
		if len(a)-i < d {
			v, err = strconv.Atoi(a[i : i+d])
		} else {
			v, err = strconv.Atoi(a[i : len(a)-i])
		}
		if err != nil {
			panic(err)
		}
		as = append(as, v)
	}
	return as
}

func BigAdd2(a string, b string) string {
	var as = qf(a)
	var bs = qf(b)
	var ds []int
	if len(as) > len(bs) {
		ds = as[len(as)-len(bs):]
		as = as[:len(bs)]
	} else {
		ds = bs[len(bs)-len(as):]
		bs = bs[:len(as)]
	}

	var v = 0
	var li = 0
	var cd = len(ds)
	for x := len(as) - 1; x >= 0; x-- {
		v = as[x] + bs[x] + li
		if v != v%1_000_000 {
			li = 1
			v = v % 1_000_000
		} else {
			li = 0
		}
		ds = append(ds, v)
	}

	if li != 0 {
		for x := cd - 1; x >= 0; x-- {
			ds[x] = ds[x] + 1
			if ds[x] != ds[x]%1_000_000 {
				ds[x] = ds[x] % 1_000_000
			} else {
				break
			}
		}
	}

	return ""
}

func BigAdd(a string, b string) string {
	if len(a) < len(b) {
		a = muchZero(len(b)-len(a)) + a
	} else {
		b = muchZero(len(a)-len(b)) + b
	}

	var x = false
	var v byte
	var vs = make([]string, 0, len(a))
	for i := len(a) - 1; i >= 0; i-- {
		if (a[i] < '0' || a[i] > '9') || (b[i] < '0' || b[i] > '9') {
			panic("err")
		}

		v, x = bigAdd(a[i], b[i], x)
		vs = append(vs, string(v))
	}

	var res = ""
	for i := len(vs) - 1; i >= 0; i-- {
		res += vs[i]
	}

	if x {
		return "1" + res
	}

	return res
}

func bigAdd(a byte, b byte, x bool) (byte, bool) {
	var d byte = 0
	if x {
		d = 1
	}
	var c = a + b - '0' + d
	if c > '9' {
		return c - '1' - '9' + '0' + '0', true
	}
	return c, false
}
