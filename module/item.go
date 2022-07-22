package module

// Item 条目，由resp中筛选
type Item map[string]interface{}

func (i Item) Valid() bool {
	return i != nil
}
