package zaperr_test

import (
	"errors"
	"regexp"

	"github.com/hori-ryota/go-testutil/testutil"
	"github.com/hori-ryota/zaperr"
	"go.uber.org/zap"
)

func Example() {
	testutil.OverwritingExampleOutputWrapper(func() {

		logger := zap.NewExample()
		defer logger.Sync()

		err := errors.New("error")

		err = zaperr.Wrap(err, "failed to execute something",
			zap.Int("foo", 1),
			zap.String("bar", "baz"),
		)

		logger.Info("example", zaperr.ToField(err))

	}, func(s []byte) []byte {
		// replace stdout for excepting the absolute path depending on execution
		// environment
		return regexp.MustCompile(`"errorVerbose":"[^"]*"`).ReplaceAll(s, []byte(`"errorVerbose":"omitted..."`))
	})
	// Output:
	// {"level":"info","msg":"example","error":{"foo":1,"bar":"baz","error":"failed to execute something: error","errorVerbose":"omitted..."}}
}
