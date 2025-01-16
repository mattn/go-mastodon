package mastodon

import (
	"encoding/json"
	"strconv"
)

type ID string

func (id *ID) UnmarshalJSON(data []byte) error {
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
//	-1 if j is greater than i.
//
// Compare can be used as an argument of [slices.SortFunc]:
//
//	slices.SortFunc([]mastodon.ID{id1, id2}, mastodon.ID.Compare)
func (i ID) Compare(j ID) int {
	var (
		ii = i.u64()
		jj = j.u64()
	)

	switch {
	case ii < jj:
		return -1
	case ii == jj:
		return 0
	case jj < ii:
		return +1
	}
	panic("impossible")
}

func (i ID) u64() uint64 {
	if i == "" {
		return 0
	}
	v, err := strconv.ParseUint(string(i), 10, 64)
	if err != nil {
		panic(err)
	}
	return v
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
