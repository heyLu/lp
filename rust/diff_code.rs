fn diff_encode(vals: &[int]) -> Box<&[int]> {
	if vals.len() < 1 {
		return box &[]
	}

	let mut encoded: &[int] = &[];
	let mut last = vals[0];

	for i in range(1, vals.len()) {
		encoded[i] = (vals[i] - last);
	}

	box encoded
}

fn main() {
	println!("diff_encode(&[]) = {}", diff_encode(&[]));

	let years = [1913i, 1981, 1960, 1920, 1980, 2023, 2807];
	println!("years: {}", years.as_slice());
	println!("encoded: {}", diff_encode(years));
}
