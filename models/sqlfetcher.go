package samasemodels

type SQLFetcher interface {
	AddJoins(joins string)
	AddFields(fields string)
	AddScanDest(dest []interface{})
}
