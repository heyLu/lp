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

struct CNF {
    num_vars: u32,
    num_clauses: u32,
    clauses: Vec<Vec<i32>>
}

fn parse_dimac(dimac: &str) -> Result<CNF, String> {
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
            println!("cnf has {} variables and {} clauses", num_vars, num_clauses)
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
        println!("{:?}", vars);
        clauses.push(vars)
    }

    let cnf = CNF { num_vars: num_vars, num_clauses: num_clauses, clauses: clauses };
    return Ok(cnf)
}

fn main() {
    let input: &mut String = &mut String::new();
    match io::stdin().read_to_string(input) {
        Ok(_) => match parse_dimac(input) {
            Ok(cnf) => { println!("ok!") }
            Err(e) => { println!("Error: {}", e) }
        },
        Err(e) => { println!("Error: {}", e) }
    }
}
