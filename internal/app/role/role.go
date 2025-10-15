package role

type Role int

const (
	User Role = iota
	Moderator
)

func (r Role) String() string {
	switch r {
	case User:
		return "user"
	case Moderator:
		return "moderator"
	default:
		return "guest"
	}
}