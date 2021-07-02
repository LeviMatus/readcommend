package api

import (
	"strings"
	"testing"

	"github.com/LeviMatus/readcommend/service/internal/driver/author"
	"github.com/LeviMatus/readcommend/service/internal/driver/author/authortest"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/driver/book/booktest"
	"github.com/LeviMatus/readcommend/service/internal/driver/era"
	"github.com/LeviMatus/readcommend/service/internal/driver/era/eratest"
	"github.com/LeviMatus/readcommend/service/internal/driver/genre"
	"github.com/LeviMatus/readcommend/service/internal/driver/genre/genretest"
	"github.com/LeviMatus/readcommend/service/internal/driver/size"
	"github.com/LeviMatus/readcommend/service/internal/driver/size/sizetest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {

	tests := map[string]struct {
		authorDriver   author.Driver
		sizeDriver     size.Driver
		genreDriver    genre.Driver
		eraDriver      era.Driver
		bookDriver     book.Driver
		expectedRoutes []string
		errAssertion   assert.ErrorAssertionFunc
		valAssertion   assert.ValueAssertionFunc
	}{
		"code contract not met - nil driver": {
			errAssertion: assert.Error,
			valAssertion: assert.Nil,
		},
		"create server": {
			authorDriver:   &authortest.DriverMock{},
			sizeDriver:     &sizetest.DriverMock{},
			genreDriver:    &genretest.DriverMock{},
			eraDriver:      &eratest.DriverMock{},
			bookDriver:     &booktest.DriverMock{},
			expectedRoutes: []string{"api"},
			errAssertion:   assert.NoError,
			valAssertion:   assert.NotNil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server, err := New(tt.authorDriver, tt.sizeDriver, tt.genreDriver, tt.eraDriver, tt.bookDriver, zap.NewNop())
			tt.errAssertion(t, err)
			tt.valAssertion(t, server)
			if server != nil {
				var foundPatterns = map[string]struct{}{}
				for _, pattern := range tt.expectedRoutes {
					for _, subroute := range server.mux.Routes() {
						if strings.Contains(subroute.Pattern, pattern) {
							foundPatterns[pattern] = struct{}{}
						}
					}
					_, found := foundPatterns[pattern]
					assert.True(t, found)
				}
			}
		})
	}
}
