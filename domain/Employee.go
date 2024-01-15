package domain

type Employee struct {
	Id         int    `json:"id"`
	EmployeeId int    `json:"employeeId"`
	Username   string `json:"username"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	AutoAssign bool   `json:"autoAssign"`
}
