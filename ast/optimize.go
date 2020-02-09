package ast

/*
type FreeVars struct {
	vars map[string]bool   // whether a var is free or bound
	refs map[string][]*AST // refs to free vars
}

func FreeVarsOf(ast AST) map[string][]*AST {
	v := &FreeVars{
		vars: map[string]bool{},
		refs: map[string][]*AST{},
	}
	ast.Visit(v, &ast)

	return v.refs
}

func (v *FreeVars) Set(id string, val *AST) {
	if free, ok := v.vars[id]; ok && !free {
		return
	}

	if _, ok := v.refs[id]; !ok {
		v.refs[id] = []*AST{}
	}

	v.vars[id] = true
	v.refs[id] = append(v.refs[id], val)
}

func (v *FreeVars) VisitApplication(e Application, aptr *AST) {
	for i, child := range e.Body {
		child.Visit(v, &e.Body[i])
	}
}

func (v *FreeVars) VisitPattern(e Pattern, aptr *AST) {
	for _, matchGroup := range e.Matches {
		for i, match := range matchGroup {
			switch M := match.(type) {
			case Identifier:
				if free, ok := v.vars[M.Value]; !ok || !free {
					v.Set(M.Value, &matchGroup[i])
				}
			default:
				match.Visit(v, &matchGroup[i])
			}
		}
	}

	for i, body := range e.Bodies {
		body.Visit(v, &e.Bodies[i])
	}
}

func (v *FreeVars) VisitIdentifier(e Identifier, aptr *AST) {
	// set identifier as free if not already bound
}

func (v *FreeVars) VisitList(e List, aptr *AST) {
	for i, child := range e.Values {
		child.Visit(v, &e.Values[i])
	}
}
func (v *FreeVars) VisitListConstructor(e ListConstructor, aptr *AST) {
	panic("TODO list constructor")
}

func (v *FreeVars) VisitWhere(e Where, aptr *AST) {
	panic("TODO where")
}

func (v *FreeVars) VisitLet(e Let, aptr *AST) {

}

func (v *FreeVars) VisitLabel(e Label, aptr *AST)   {}
func (v *FreeVars) VisitString(e String, aptr *AST) {}
func (v *FreeVars) VisitNumber(e Number, aptr *AST) {}
*/
