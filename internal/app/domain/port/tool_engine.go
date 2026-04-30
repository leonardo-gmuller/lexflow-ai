package port

type ToolEngine interface {
	Execute(input string, args map[string]any) (string, error)
	GetTools() []map[string]any
}
