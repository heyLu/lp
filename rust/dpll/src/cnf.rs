//! The cnf library provides facilities to work with CNF formulas.
//!
//! `CNF` stands for conjunctive normal form, i.e. a logic formula that
//! consists of a conjunction of disjunctions.

pub struct CNF {
	pub num_vars: u32,
		 pub num_clauses: u32,
		 pub clauses: Vec<Vec<i32>>
}

/// Parses a CNF formula in DIMAC format.
///
/// # Examples
///
/// A simple example for a formula in DIMAC format is below:
///
/// ```text
/// p cnf 5 3
/// 1 2 3 0
/// -2 3 4 0
/// 4 5 0
/// ```
///
/// The above represents a formula with 5 variables and 3 clauses.  The first
/// line specifies this.  Each following line represents a clause with possibly
/// negated literals, terminated by 0 and a newline.
pub fn parse_dimac(dimac: &str) -> Result<CNF, String> {
	let mut lines = dimac.lines();
	let mut num_vars;
	let mut num_clauses;

	match lines.next() {
		None => { return Err("expected cnf description".to_string()) }
		Some(line) => {
			let desc: Vec<&str> = line.split(" ").collect();
			if desc.len() != 4 || desc[0] != "p" || desc[1] != "cnf" {
				return Err("cnf description must be of the form 'p cnf <num vars> <num clauses>'".to_string())
			}
			match desc[2].parse::<u32>() {
				Ok(n) => { num_vars = n }
				Err(e) => { return Err(format!("<num vars> must be a positive integer: {}", e)) }
			}

			match desc[3].parse::<u32>() {
				Ok(n) => { num_clauses = n }
				Err(e) => { return Err(format!("<num clauses> must be a positive integer: {}", e)) }
			}
		}
	}

	let clause_lines: Vec<&str> = lines.collect();
	if clause_lines.len() as u32 != num_clauses {
		return Err(format!("Wrong number of clauses: Expected {}, but got {}", num_clauses, clause_lines.len()))
	}

	let mut clauses: Vec<Vec<i32>> = Vec::with_capacity(num_clauses as usize);
	for clause_line in clause_lines {
		let mut vars: Vec<i32> = clause_line.split(" ").map(|x| x.parse::<i32>().unwrap()).collect();
		if vars.is_empty() {
			return Err("empty clause".to_string())
		}
		if vars[vars.len()-1] != 0 {
			return Err("clause must be terminated with 0".to_string())
		}
		let l = vars.len();
		vars.truncate(l - 1);
		clauses.push(vars)
	}

	let cnf = CNF { num_vars: num_vars, num_clauses: num_clauses, clauses: clauses };
	return Ok(cnf)
}
