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

/// A satisfied clause is a clause where at least one atom is true.
fn is_clause_satisfied(vars: BoundVars, clause: Clause) -> bool {
    for v in clause {
        if vars.contains(&v) {
            return true
        }
    }
    
    return false
}

#[test]
fn test_is_clause_satisfied() {
    assert!(is_clause_satisfied(from_vec(vec!(1)), vec!(1)));
    assert!(!is_clause_satisfied(empty_vars(), vec!(1)));
}

/// A conflict clause is a clause whose atoms each are false.
fn is_clause_conflict(vars: BoundVars, clause: Clause) -> bool {
    for v in clause {
        if vars.contains(&v) {
            return false
        }
    }

    return true
}

#[test]
fn test_is_clause_conflict() {
    assert!(is_clause_conflict(empty_vars(), vec!(1)));
    assert!(is_clause_conflict(empty_vars(), vec!(1, 2, 3)));
    assert!(is_clause_conflict(from_vec(vec!(4)), vec!(1, 2, 3)));
    assert!(!is_clause_conflict(from_vec(vec!(1)), vec!(1)));
    assert!(!is_clause_conflict(from_vec(vec!(2)), vec!(1, 2)));
}
