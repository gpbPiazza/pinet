package pihttp

// key value pair
// headers can have repetitive key names, if so,
// this values will me concatenated into the same key
type header map[string][]string
