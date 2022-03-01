package task

const (
	EntityType = "task"
)

var validPropertySysnames = []string{}

// Task ...
type Task struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// New func is a constructor for the Task
func New() *Task {
	return &Task{}
}

func (e *Task) GetValidPropertySysnames() []string {
	return validPropertySysnames
}

func (e *Task) GetMapNameSysname() map[string]string {
	return map[string]string{
		//"Email":     PropertySysnameEmail,
		//"AccountID": PropertySysnameAccountID,
	}
}
