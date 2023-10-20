/**
 * @Time: 2023/10/20 16:40
 * @Author: jzechen
 * @File: errcode.go
 * @Software: GoLand collector
 */

package errcode

// ErrCode define error code
type ErrCode int

//go:generate stringer -type ErrCode -linecomment
const (
	// Common Errors
	OK          ErrCode = iota // success
	DBErr                      // database operation error
	JsonMarshal                // json marshall or unmarshall error
	ServerErr                  // server error
	NotFound                   // not found

	// TODO: add error code here
)

func (ec ErrCode) Error() string {
	return ec.String()
}
