package nc

const documentation = []struct {
	Endpoint    string
	Request     string
	Response    string
	Description string
}{
	{"ari:applications:all", "ignored", "[]string", "Get all applications"},
	{"ari:applications:data:>", "ignored", "ari.ApplicationData", "Get the data for the given applications"},
	{"ari:applications:subscribe:>", "string", "ignored", "Subscribe the app to the event source"},
	{"ari:applications:unsubscribe:>", "string", "ignored", "Unsubscribe the app from the event source"},
}
