package actions

type Action interface {
	//all an action can do
	Execute() (interface{}, error)
}