/*
 * Reads a CNF formula in DIMAC format on stdin, parses it and prints
 * the result.
 *
 * Many things (support for comments, empty lines and sane code in
 * general) are still missing.
 *
 * You have been warned ...
 */

extern crate solve;

use std::io;
use std::io::Read;

use solve::cnf;

fn main() {
	let input: &mut String = &mut String::new();
	match io::stdin().read_to_string(input) {
		Ok(_) => match cnf::parse_dimac(input) {
			Ok(cnf) => {
				println!("cnf has {} variables and {} clauses", cnf.num_vars, cnf.num_clauses);
				for clause in cnf.clauses {
					println!("{:?}", clause);
				}
			}
			Err(e) => { println!("Error: {}", e) }
		},
			Err(e) => { println!("Error: {}", e) }
	}
}
