package stuc

// Dinsight Generate dinsight software authorization file
func Dinsight(app, ver, awd string, typ, exp int, start int64, per string) ([]byte, error) {
	return New(app, ver, awd, typ, exp, start, per)
}
