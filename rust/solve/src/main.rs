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
use std::thread;

use solve::cnf;
use solve::dpll;

fn main() {
    let input: &mut String = &mut String::new();
    match io::stdin().read_to_string(input) {
	Ok(_) => match cnf::parse_dimac(input) {
	    Ok(cnf) => {
		println!("cnf has {} variables and {} clauses", cnf.num_vars, cnf.num_clauses);
		for clause in cnf.clauses.clone() {
		    println!("{:?}", clause);
		}

                thread::Builder::new().name("solver".to_string()).stack_size(100 * 1024 * 1024).spawn(move || {
                    match dpll::dpll(cnf.clauses) {
                        Some(bindings) => println!("satisfiable: {:?}", bindings),
                        None => println!("not satisfiable")
                    }
                }).unwrap().join();
	    }
	    Err(e) => { println!("Error: {}", e) }
	},
	Err(e) => { println!("Error: {}", e) }
    }
}
