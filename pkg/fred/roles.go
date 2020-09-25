package fred

// Role is a wrapper for a set of permissions/methods.
type Role string

// These are all roles that we have in FReD.
const (
	ReadKeygroup       Role = "R"
	WriteKeygroup      Role = "W"
	ConfigureReplica   Role = "C"
	ConfigureTrigger   Role = "T"
	ConfigureKeygroups Role = "K"
)

var (
	permissions = map[Role]map[Method]struct{}{
		ReadKeygroup: {
			Read: {},
		},
		WriteKeygroup: {
			Update: {},
			Delete: {},
		},
		ConfigureReplica: {
			AddReplica:    {},
			GetReplica:    {},
			RemoveReplica: {},
		},
		ConfigureTrigger: {
			GetTrigger:    {},
			AddTrigger:    {},
			RemoveTrigger: {},
		},
		ConfigureKeygroups: {
			DeleteKeygroup: {},
			AddUser:        {},
			RemoveUser:     {},
		},
	}
)
