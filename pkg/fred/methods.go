package fred

// Method is a representation of all methods a client could perform on a keygroup.
type Method string

// These are all methods that clients can perform on FReD, implemented as constants for easier use.
const (
	CreateKeygroup Method = "CreateKeygroup"
	DeleteKeygroup Method = "DeleteKeygroup"
	Read           Method = "Read"
	Update         Method = "Update"
	Delete         Method = "Delete"
	AddReplica     Method = "AddReplica"
	GetReplica     Method = "GetReplica"
	RemoveReplica  Method = "RemoveReplica"
	GetAllReplica  Method = "GetAllReplica"
	GetTrigger     Method = "GetTrigger"
	AddTrigger     Method = "AddTrigger"
	RemoveTrigger  Method = "RemoveTrigger"
	AddUser        Method = "AddUser"
	RemoveUser     Method = "RemoveUser"
)
