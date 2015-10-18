package epm

import (
	"fmt"

	// "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/modules/types"
)

// TODO: abstract away the mechanism here so the particular
// chain module can deal with state diffs and we don't need silly State
// will be faster too

func (e *EPM) CurrentState() State { //map[string]string{
	if e.chain == nil {
		return State{}
	}
	return State{}
}

func (e *EPM) newDiffSched(i int) {
	if e.diffSched[i] == nil {
		e.diffSched[i] = []string{}
	}
}

func (e *EPM) checkTakeStateDiff(i int) {
	if _, ok := e.diffSched[i]; !ok {
		return
	}
	e.Commit()
	names := e.diffSched[i]
	for _, name := range names {
		if _, ok := e.states[name]; !ok {
			// store state
			e.states[name] = e.CurrentState()
		} else {
			// take diff
			e.Commit()
			PrintDiff(name, e.states[name], e.CurrentState())
		}
	}
}

func StorageDiff(pre, post State) State { //map[string]string) map[string]map[string]string{
	diff := State{make(map[string]*Storage), []string{}}
	// for each account in post, compare all elements.
	for _, addr := range post.Order {
		acct := post.State[addr]
		diff.State[addr] = &Storage{make(map[string]string), []string{}}
		diff.Order = append(diff.Order, addr)
		acct2, ok := pre.State[addr]
		if !ok {
			// if this account didnt exist in pre
			diff.State[addr] = acct
			continue
		}
		// for each storage in the post acct, check for diff in 2.
		for _, k := range acct.Order {
			v := acct.Storage[k]
			v2, ok := acct2.Storage[k]
			// if its not in the pre-state or its different, add to diff
			if !ok || v2 != v {
				diff.State[addr].Storage[k] = v
				st := diff.State[addr]
				st.Order = append(diff.State[addr].Order, k)
				diff.State[addr] = st
			}
		}
	}
	return diff
}

func PrettyPrintAcctDiff(dif State) string { //map[string]string) string{
	result := ""
	for _, addr := range dif.Order {
		acct := dif.State[addr]
		if len(acct.Order) == 0 {
			continue
		}
		result += addr + ":\n"
		for _, store := range acct.Order {
			v := acct.Storage[store]
			val := v
			result += "\t" + store + ": " + val + "\n"
		}
	}
	return result
}

func PrintDiff(name string, pre, post State) { //map[string]string) {
	/*
	   fmt.Println("pre")
	   fmt.Println(PrettyPrintAcctDiff(pre))
	   fmt.Println("\n\n")
	   fmt.Println("post")
	   fmt.Println(PrettyPrintAcctDiff(post))
	   fmt.Println("\n\n")
	*/
	fmt.Println("diff:", name)
	diff := StorageDiff(pre, post)
	fmt.Println(PrettyPrintAcctDiff(diff))
}
