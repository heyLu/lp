use std::collections::BTreeSet;

type Var = i32;
type BoundVars = BTreeSet<Var>;
type Clause = Vec<Var>;

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
fn is_clause_satisfied(vars: &BoundVars, clause: Clause) -> bool {
    for v in clause {
        if is_true(vars, v) {
            return true
        }
    }
    
    return false
}

#[test]
fn test_is_clause_satisfied() {
    assert!(is_clause_satisfied(&from_vec(vec!(1)), vec!(1)));
    assert!(!is_clause_satisfied(&empty_vars(), vec!(1)));
}

/// A conflict clause is a clause whose atoms each are false.
fn is_clause_conflict(vars: &BoundVars, clause: Clause) -> bool {
    for v in clause {
        if is_true(vars, v) {
            return false
        }
    }

    return true
}

#[test]
fn test_is_clause_conflict() {
    assert!(is_clause_conflict(&empty_vars(), vec!(1)));
    assert!(is_clause_conflict(&empty_vars(), vec!(1, 2, 3)));
    assert!(is_clause_conflict(&from_vec(vec!(4)), vec!(1, 2, 3)));
    assert!(!is_clause_conflict(&from_vec(vec!(1)), vec!(1)));
    assert!(!is_clause_conflict(&from_vec(vec!(2)), vec!(1, 2)));
}
