use std::collections::BTreeSet;

pub type Var = i32;
pub type BoundVars = BTreeSet<Var>;
pub type Clause = Vec<Var>;

fn empty_vars() -> BoundVars {
    BTreeSet::new()
}

fn from_vec(v: Vec<Var>) -> BoundVars {
    let mut vars = empty_vars();
    vars.extend(v);
    vars
}

fn is_true(vars: &BoundVars, var: Var) -> bool {
    let t = vars.contains(&var);
    let nt = vars.contains(&-var);
    if t || nt {
        t || !nt
    } else {
        false
    }
}

#[test]
fn test_is_true() {
    let vars = &from_vec(vec!(1, -2));
    assert!(is_true(vars, 1));
    assert!(is_true(vars, -2));
    assert!(!is_true(vars, -1));
    assert!(!is_true(vars, 2));
    assert!(!is_true(vars, 3));
    assert!(!is_true(vars, -3));
}

fn is_false(vars: &BoundVars, var: Var) -> bool {
    let t = vars.contains(&var);
    let nt = vars.contains(&-var);
    if t || nt {
        !t || nt
    } else {
        false
    }
}

#[test]
fn test_is_false() {
    let vars = &from_vec(vec!(1, -2));
    assert!(is_false(vars, -1));
    assert!(is_false(vars, 2));
    assert!(!is_false(vars, 1));
    assert!(!is_false(vars, -2));
    assert!(!is_false(vars, 3));
    assert!(!is_false(vars, -3));
}

fn is_unknown(vars: &BoundVars, var: Var) -> bool {
    !vars.contains(&var) && !vars.contains(&-var)
}

#[test]
fn test_is_unknown() {
    let vars = &from_vec(vec!(1, -2));
    assert!(is_unknown(vars, 3));
    assert!(is_unknown(vars, -3));
    assert!(is_unknown(vars, 4));
    assert!(is_unknown(vars, -4));
    assert!(!is_unknown(vars, 1));
    assert!(!is_unknown(vars, -1));
    assert!(!is_unknown(vars, 2));
    assert!(!is_unknown(vars, -2));
}

/// A satisfied clause is a clause where at least one atom is true.
fn is_clause_satisfied(vars: &BoundVars, clause: &Clause) -> bool {
    clause.iter().any(|v| is_true(vars, *v))
}

#[test]
fn test_is_clause_satisfied() {
    assert!(is_clause_satisfied(&from_vec(vec!(1)), &vec!(1)));
    assert!(!is_clause_satisfied(&empty_vars(), &vec!(1)));
}

/// A conflict clause is a clause whose atoms each are false.
fn is_clause_conflict(vars: &BoundVars, clause: &Clause) -> bool {
    clause.iter().all(|v| is_false(vars, *v))
}

#[test]
fn test_is_clause_conflict() {
    assert!(is_clause_conflict(&from_vec(vec!(-1)), &vec!(1)));
    assert!(is_clause_conflict(&from_vec(vec!(-1, -2, -3)), &vec!(1, 2, 3)));
    assert!(is_clause_conflict(&from_vec(vec!(-1, -2, -3, 4)), &vec!(1, 2, 3)));
    assert!(!is_clause_conflict(&from_vec(vec!(1)), &vec!(1)));
    assert!(!is_clause_conflict(&from_vec(vec!(2)), &vec!(1, 2)));
}

/// A unit clause is a clause where one atom is unknown and all
/// others are false.
fn is_clause_unit(vars: &BoundVars, clause: &Clause) -> bool {
    let mut unknowns = 0;
    
    for &v in clause {
        if is_unknown(vars, v) {
            unknowns += 1
        }

        if unknowns > 1 || is_true(vars, v) {
            return false
        }
    }

    return unknowns == 1
}

#[test]
fn test_is_clause_unit() {
    let vars = &from_vec(vec!(1, 2));
    assert!(is_clause_unit(vars, &vec!(3)));
    assert!(is_clause_unit(vars, &vec!(-1, 3)));
    assert!(is_clause_unit(vars, &vec!(-1, -2, 3)));
    assert!(!is_clause_unit(vars, &vec!(1, 3)));
    assert!(!is_clause_unit(vars, &vec!(1, 2, 3)));
    assert!(!is_clause_unit(vars, &vec!(1, 2)));
}

pub fn dpll(clauses: Vec<Clause>) -> Option<BoundVars> {
    let stack: &mut Vec<(Var, BoundVars)> = &mut Vec::new();
    let mut vars = &mut empty_vars();
    
    loop {
        if clauses.iter().all(|c| is_clause_satisfied(&vars, c)) { // all clauses satisfied, success
            //println!("satisfied");
            return Some(vars.clone())
        } else if clauses.iter().any(|c| is_clause_conflict(&vars, c)) { // a conflict exists, backtrack
            match stack.pop() {
                None => return None, // nothing to backtrack, no solution found
                Some((v, b)) => {
                    //println!("backtrack! {}", v);
                    vars.clone_from(&b);
                    vars.insert(v);
                }
            }
        } else if clauses.iter().any(|c| is_clause_unit(&vars, c)) { // a unit clause exists, propagate
            let cs = clauses.clone();
            let clause = cs.iter().find(|&c| is_clause_unit(&vars, c)).unwrap();
            let unknown = *clause.iter().find(|&v| is_unknown(&vars, *v)).unwrap();
            //println!("propagate {} from {:?}", unknown, clause);
            vars.insert(unknown);
        } else { // none of the above, decide (guess) an unknown variable
            let mut unknown: Var = 0;
            for c in clauses.clone() {
                for v in c {
                    if is_unknown(vars, v) {
                        unknown = v;
                        break;
                    }
                }
                
                if unknown != 0 {
                    break;
                }
            }
            assert!(unknown != 0);
            stack.push((unknown, vars.clone()));
            //println!("decide {}", -unknown);
            vars.insert(-unknown);
        }
    }
}

#[test]
fn test_dpll_trivial() {
    let clauses = vec!(vec!(1, 2), vec!(2));
    assert!(dpll(clauses).is_some());
}
