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
