package group

// Group of people, each of which have a balance in the group.
type Group struct {
	ID      int
	Name    string
	Members map[string]float32
}

// AddMember to the given group.
func (g *Group) AddMember(member string) error {
	if g.Members == nil {
		g.Members = make(map[string]float32)
	}

	if _, prs := g.Members[member]; prs {
		return &ExistingMembersError{g.ID, []string{member}}
	}

	g.Members[member] = 0.0

	return nil
}

// AddMembers to the given group.
func (g *Group) AddMembers(members []string) error {
	if g.Members == nil {
		g.Members = make(map[string]float32)
	}

	var exs []string

	for _, m := range members {
		if _, prs := g.Members[m]; prs {
			exs = append(exs, m)
			continue
		}

		g.Members[m] = 0.0
	}

	if len(exs) != 0 {
		return &ExistingMembersError{g.ID, exs}
	}

	return nil
}
