package params

import (
	"net/http"
	"strconv"

	"github.com/LeviMatus/readcommend/service/pkg/util"
	"github.com/pkg/errors"
)

func String(r *http.Request, n string) *string {
	val := r.URL.Query().Get(n)
	return util.StringPtr(val)
}

func Int16(r *http.Request, n string) (*int16, error) {
	param := r.URL.Query().Get(n)
	val, err := strconv.ParseInt(param, 10, 16)
	if errors.Is(err, strconv.ErrSyntax) && param == "" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	out := int16(val)
	return &out, err
}

func Uint64(r *http.Request, n string) (*uint64, error) {
	param := r.URL.Query().Get(n)
	val, err := strconv.ParseUint(param, 10, 64)
	if errors.Is(err, strconv.ErrSyntax) && param == "" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &val, err
}

func Int16Slice(r *http.Request, n string) ([]int16, error) {
	// The following could be broken up into its own StringSlice routine.
	// Its not needed right now, but if we wanted to use []int32, for example, we
	// could just make a StringSlice and call it from Int16Slice and Int32Slice in the future.
	vals, ok := r.URL.Query()[n]
	if !ok {
		return nil, nil
	}
	out := make([]int16, len(vals))
	for i, val := range vals {
		intVal, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return nil, err
		}
		out[i] = int16(intVal)
	}
	return out, nil
}
