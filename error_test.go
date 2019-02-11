package zaperr_test

import (
	commonerr "errors"
	"fmt"
	"testing"

	"github.com/hori-ryota/zaperr"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestNew(t *testing.T) {
	var err error
	err = zaperr.New("error",
		zap.Int("1", 1),
		zap.Int("2", 2),
	)

	t.Run("Error", func(t *testing.T) {
		assert.EqualError(t, err, "error")
		assert.Contains(t, fmt.Sprintf("%+v", err), "error\n")
	})
	t.Run("log", func(t *testing.T) {
		core, observed := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		logger.Debug("field", zaperr.ToField(err))
		logged := observed.All()[0].ContextMap()
		errElem := logged["error"].(map[string]interface{})
		assert.EqualValues(t, 1, errElem["1"])
		assert.EqualValues(t, 2, errElem["2"])
		assert.Equal(t, "error", errElem["error"])
	})
}

func TestErrorf(t *testing.T) {
	var err error
	err = zaperr.Errorf("error %s", "value")
	assert.EqualError(t, err, "error value")
	assert.Contains(t, fmt.Sprintf("%+v", err), "error value\n")
}

func TestWrap(t *testing.T) {

	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, zaperr.Wrap(nil, "wrap"))
	})

	t.Run("Wrap", func(t *testing.T) {
		var err error
		err = zaperr.New("error")
		err = zaperr.Wrap(err, "wrap")
		assert.EqualError(t, err, "wrap: error")
		assert.Contains(t, fmt.Sprintf("%+v", err), "error\n")
		assert.Contains(t, fmt.Sprintf("%+v", err), "wrap\n")
		err = zaperr.Wrap(err, "wrap2")
		assert.EqualError(t, err, "wrap2: wrap: error")
	})

	t.Run("Wrap: with fields", func(t *testing.T) {
		var err error
		err = zaperr.New("error")
		err = zaperr.Wrap(err, "wrap")
		assert.EqualError(t, err, "wrap: error")
		assert.Contains(t, fmt.Sprintf("%+v", err), "error\n")
		assert.Contains(t, fmt.Sprintf("%+v", err), "wrap\n")
	})

	t.Run("Wrap: with fields: nested", func(t *testing.T) {
		core, observed := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		var err error
		err = zaperr.New("error")
		err = zaperr.Wrap(
			err,
			"wrap",
			zap.Int("1", 1),
			zap.Int("2", 2),
		)
		assert.EqualError(t, err, "wrap: error")
		err = zaperr.Wrap(
			err,
			"wrap2",
			zap.Int("3", 3),
			zap.Int("4", 4),
		)
		assert.EqualError(t, err, "wrap2: wrap: error")

		logger.Debug("field", zaperr.ToField(err))
		logged := observed.All()[0].ContextMap()
		errElem := logged["error"].(map[string]interface{})
		assert.EqualValues(t, 3, errElem["3"])
		assert.EqualValues(t, 4, errElem["4"])
		assert.EqualValues(t, "wrap2", errElem["message"])
		nestedElem := errElem["error"].(map[string]interface{})
		assert.EqualValues(t, 1, nestedElem["1"])
		assert.EqualValues(t, 2, nestedElem["2"])
		assert.Equal(t, "wrap: error", nestedElem["error"])
	})
}

func TestCause(t *testing.T) {
	var err error
	err = zaperr.New("error")
	err = zaperr.Wrap(err, "wrap")
	err = zaperr.Wrap(err, "wrap")
	assert.EqualError(t, errors.Cause(err), "error")
}

func TestToField(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		core, observed := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		logger.Debug("field", zaperr.ToField(nil))
		assert.Empty(t, observed.All()[0].ContextMap())
	})
	t.Run("WithFields", func(t *testing.T) {
		core, observed := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		var err error
		err = zaperr.New("error")
		err = zaperr.WithFields(
			err,
			zap.Int("1", 1),
			zap.Int("2", 2),
		)
		logger.Debug("field", zaperr.ToField(err))
		logged := observed.All()[0].ContextMap()
		errElem := logged["error"].(map[string]interface{})
		assert.EqualValues(t, 1, errElem["1"])
		assert.EqualValues(t, 2, errElem["2"])
		assert.Equal(t, "error", errElem["error"])
	})
	t.Run("WithFields: append", func(t *testing.T) {
		core, observed := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		var err error
		err = zaperr.New("error")
		err = zaperr.WithFields(
			err,
			zap.Int("1", 1),
			zap.Int("2", 2),
		)
		err = zaperr.WithFields(
			err,
			zap.Int("3", 3),
			zap.Int("4", 4),
		)
		logger.Debug("field", zaperr.ToField(err))
		logged := observed.All()[0].ContextMap()
		errElem := logged["error"].(map[string]interface{})
		assert.EqualValues(t, 1, errElem["1"])
		assert.EqualValues(t, 2, errElem["2"])
		assert.EqualValues(t, 3, errElem["3"])
		assert.EqualValues(t, 4, errElem["4"])
		assert.Equal(t, "error", errElem["error"])
	})
	t.Run("WithFields: wrapped", func(t *testing.T) {
		core, observed := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		var err error
		err = zaperr.New("error")
		err = zaperr.WithFields(
			err,
			zap.Int("1", 1),
			zap.Int("2", 2),
		)
		err = zaperr.Wrap(err, "wrap")
		logger.Debug("field", zaperr.ToField(err))
		logged := observed.All()[0].ContextMap()
		errElem := logged["error"].(map[string]interface{})
		assert.EqualValues(t, 1, errElem["1"])
		assert.EqualValues(t, 2, errElem["2"])
		assert.Equal(t, "wrap: error", errElem["error"])
	})
	t.Run("without fields", func(t *testing.T) {
		core, observed := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		var err error
		err = errors.New("error")
		logger.Debug("field", zaperr.ToField(err))
		logged := observed.All()[0].ContextMap()
		assert.Equal(t, "error", logged["error"])
	})
	t.Run("without fields: wrapped", func(t *testing.T) {
		core, observed := observer.New(zap.DebugLevel)
		logger := zap.New(core)

		var err error
		err = errors.New("error")
		err = errors.Wrap(err, "wrap")
		logger.Debug("field", zaperr.ToField(err))
		logged := observed.All()[0].ContextMap()
		assert.Equal(t, "wrap: error", logged["error"])
	})
}

func TestToNamedField(t *testing.T) {
	for _, tt := range []struct {
		name   string
		source error
	}{
		{
			name:   "zaperr",
			source: zaperr.New("error"),
		},
		{
			name:   "github.com/pkg/errors",
			source: errors.New("error"),
		},
		{
			name:   "common error",
			source: commonerr.New("error"),
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			core, observed := observer.New(zap.DebugLevel)
			logger := zap.New(core)

			logger.Debug("field", zaperr.ToNamedField("name", tt.source))
			logged := observed.All()[0].ContextMap()
			assert.Empty(t, logged["error"])
			assert.NotEmpty(t, logged["name"])
		})
	}
}
