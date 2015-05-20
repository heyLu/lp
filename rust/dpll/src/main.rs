/*
 * Terrible code for reading CNF formulas in DIMAC format.
 *
 * Don't read this, it's horrible.  Although I'm hoping to clean it up
 * at some point, I'm not sure when (if ever) that will be.
 *
 * You have been warned...
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
