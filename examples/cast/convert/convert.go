package convert

// Placeholder represents a custom type alias for a boolean.
//
// A constant boolean can be assigned to a Placeholder.
// var placeholder Placeholder
// placeholder = true
//
// To assign a Placeholder variable to a boolean variable, type conversion must occur.
// var b boolean
// boolean = bool(placeholder)
//
// To assign a boolean variable to a Placeholder variable, type conversion must occur.
// placeholder = Placeholder(boolean)
type Placeholder interface{}
