package mastodon

import (
	"encoding/json"
	"strconv"
	"strings"
)

type ID string

func (id *ID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*id = ""
		return nil
	}
	if len(data) > 0 && data[0] == '"' && data[len(data)-1] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*id = ID(s)
		return nil
	}
	var n int64
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*id = ID(strconv.FormatInt(n, 10))
	return nil
}

// Compare compares the Mastodon IDs i and j.
// Compare returns:
//
//	-1 if i is less than j,
//	 0 if i equals j,
//	+1 if i is greater than j.
//
// Shorter IDs sort before longer ones and IDs of the same length are
// compared lexicographically. This matches the numeric order of the
// integer IDs used by Mastodon while also supporting non-numeric IDs
// (e.g. ULIDs) used by other implementations.
//
// Compare can be used as an argument of [slices.SortFunc]:
//
//	slices.SortFunc([]mastodon.ID{id1, id2}, mastodon.ID.Compare)
func (i ID) Compare(j ID) int {
	if len(i) != len(j) {
		if len(i) < len(j) {
			return -1
		}
		return +1
	}
	return strings.Compare(string(i), string(j))
}

type Sbool bool

func (s *Sbool) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' && data[len(data)-1] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		b, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		*s = Sbool(b)
		return nil
	}
	var b bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	*s = Sbool(b)
	return nil
}
