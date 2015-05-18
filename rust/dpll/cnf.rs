/*
 * Terrible code for reading CNF formulas in DIMAC format.
 *
 * Don't read this, it's horrible.  Although I'm hoping to clean it up
 * at some point, I'm not sure when (if ever) that will be.
 *
 * You have been warned...
 */

use std::io;
use std::io::Read;

fn parse_dimac(dimac: &str) {
    let mut lines = dimac.lines();
    let mut num_vars = 0;
    let mut num_clauses = 0;
    
    match lines.next() {
        None => { println!("Error: expected cnf description"); return }
        Some(line) => {
            let desc: Vec<&str> = line.split(" ").collect();
            if desc.len() != 4 || desc[0] != "p" || desc[1] != "cnf" {
                println!("Error: cnf description must be of the form 'p cnf <num vars> <num clauses>'");
                return;
            }
            match desc[2].parse::<u32>() {
                Ok(n) => { num_vars = n }
                Err(e) => { println!("Error: <num vars> must be a positive integer: {}", e); return; }
            }
            
            match desc[3].parse::<u32>() {
                Ok(n) => { num_clauses = n }
                Err(e) => { println!("Error: <num clauses> must be a positive integer: {}", e); return; }
            }
            println!("cnf has {} variables and {} clauses", num_vars, num_clauses)
        }
    }

    let clause_lines: Vec<&str> = lines.collect();
    if clause_lines.len() as u32 != num_clauses {
        println!("Error: Wrong number of clauses: Expected {}, but got {}", num_clauses, clause_lines.len());
        return
    }
}

fn main() {
    let input: &mut String = &mut String::new();
    io::stdin().read_to_string(input);
    parse_dimac(input)
}
