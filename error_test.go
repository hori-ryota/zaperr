package zaperr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ fieldser = zaperr{}
var _ fieldsAppender = &zaperr{}
var _ zapcore.ObjectMarshaler = zaperr{}

func Test_zaperr_Error(t *testing.T) {
	assert.Equal(t, "error", zaperr{err: errors.New("error")}.Error())
}

func Test_zaperr_AppendFields(t *testing.T) {
	err := zaperr{
		err: errors.New("error"),
		fields: []zapcore.Field{
			zap.Int("1", 1),
		},
	}
	err.AppendFields(zap.Int("2", 2))
	assert.Equal(t,
		[]zapcore.Field{
			zap.Int("1", 1),
			zap.Int("2", 2),
			zap.Error(errors.New("error")),
		},
		err.Fields(),
	)
}

func TestAppendFields(t *testing.T) {
	t.Run("nil: return nil", func(t *testing.T) {
		assert.Nil(t, AppendFields(nil))
	})
	t.Run("not zaperr: return new zaperr", func(t *testing.T) {
		err := AppendFields(errors.New("error"), zap.Int("1", 1))
		assert.Equal(t,
			[]zapcore.Field{
				zap.Int("1", 1),
				zap.Error(errors.New("error")),
			},
			Fields(err),
		)
	})
	t.Run("zaperr", func(t *testing.T) {
		var err error
		err = &zaperr{
			err: errors.New("error"),
			fields: []zapcore.Field{
				zap.Int("1", 1),
			},
		}
		err = AppendFields(err, zap.Int("2", 2))
		t.Log(fmt.Sprintf("%#v", err))
		assert.Equal(t,
			[]zapcore.Field{
				zap.Int("1", 1),
				zap.Int("2", 2),
				zap.Error(errors.New("error")),
			},
			Fields(err),
		)
	})
}

func TestToField(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   error
		want zapcore.FieldType
	}{
		{
			name: "zaperr: Object",
			in:   zaperr{},
			want: zapcore.ObjectMarshalerType,
		},
		{
			name: "common error: not object",
			in:   errors.New("error"),
			want: zapcore.ErrorType,
		},
		{
			name: "nil: skip",
			in:   nil,
			want: zapcore.SkipType,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := ToField("name", tt.in)
			assert.Equal(t, tt.want, f.Type)
		})
	}
}
