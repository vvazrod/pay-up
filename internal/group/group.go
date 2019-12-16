package group

// Group of people, each of which have a balance in the group.
type Group struct {
	ID      int
	Name    string
	Members map[string]float32
}

// AddMember to a group.
func (g *Group) AddMember(member string) error {
	if g.Members == nil {
		g.Members = make(map[string]float32)
	}

	if _, prs := g.Members[member]; prs {
		return &ExistingMemberError{g.ID, member}
	}

	g.Members[member] = 0.0

	return nil
}

// DeleteMember from a group.
func (g *Group) DeleteMember(member string) error {
	if balance, prs := g.Members[member]; !prs {
		return &MemberNotFoundError{g.ID, member}
	} else if balance != 0.0 {
		return &DeletingBalanceError{g.ID, member, balance}
	}

	delete(g.Members, member)

	return nil
}
