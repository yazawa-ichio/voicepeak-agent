package voicepeakagent

type Narrator string

const (
	Female1     = Narrator("Female 1")
	Female2     = Narrator("Female 2")
	Female3     = Narrator("Female 3")
	Male1       = Narrator("Male 1")
	Male2       = Narrator("Male 2")
	Male3       = Narrator("Male 3")
	FemaleChild = Narrator("Female Child")
)

func (n *Narrator) String() string {
	return string(*n)
}
