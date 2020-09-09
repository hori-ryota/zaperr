package zaperr

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(message string, fields ...zap.Field) error {
	return zaperr{
		source: errors.New(message),
		fields: fields,
	}

}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func Wrap(err error, message string, fields ...zap.Field) error {
	if err == nil {
		return nil
	}
	if err, ok := err.(zaperr); ok {
		if len(fields) == 0 {
			err.source = errors.Wrap(err.source, message)
			return err
		}
		if len(err.fields) == 0 {
			return WithFields(errors.Wrap(err.source, message), fields...)
		}
		return zaperr{
			source:  err,
			fields:  fields,
			message: message,
		}
	}
	return WithFields(errors.Wrap(err, message), fields...)
}

type zaperr struct {
	source  error
	fields  []zap.Field
	message string
}

func (e zaperr) Error() string {
	if e.message != "" {
		return e.message + ": " + e.source.Error()
	}
	return e.source.Error()
}

func (e zaperr) Cause() error {
	return e.source
}

func (e zaperr) Unwrap() error {
	return e.Cause()
}

func WithFields(err error, fields ...zap.Field) error {
	if err == nil {
		return nil
	}
	if err, ok := err.(zaperr); ok {
		err.fields = append(err.fields, fields...)
		return err
	}
	return zaperr{
		source: err,
		fields: fields,
	}
}

func ToField(err error) zap.Field {
	return ToNamedField("error", err)
}

func ToNamedField(key string, err error) zap.Field {
	if err == nil {
		return zap.Skip()
	}
	if e, ok := err.(zapcore.ObjectMarshaler); ok {
		return zap.Object(key, e)
	}
	return zap.NamedError(key, err)
}

func (e zaperr) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, f := range e.fields {
		f.AddTo(enc)
	}
	if e.message != "" {
		zap.String("message", e.message).AddTo(enc)
	}
	ToField(e.source).AddTo(enc)
	return nil
}

func (e zaperr) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", e.source)
			if e.message != "" {
				fmt.Fprint(s, e.message)
			}
			return
		}
		fallthrough
	case 's', 'q':
		fmt.Fprint(s, e.Error())
	}
}

type Wrapper struct {
	fields []zap.Field
}

func WrapperWith(fields ...zap.Field) Wrapper {
	return Wrapper{fields: fields}
}

func (w Wrapper) New(message string, fields ...zap.Field) error {
	return New(message, append(fields, w.fields...)...)
}

func (w Wrapper) Errorf(format string, args ...interface{}) error {
	return WithFields(Errorf(format, args...), w.fields...)
}

func (w Wrapper) Wrap(err error, message string, fields ...zap.Field) error {
	return Wrap(err, message, append(fields, w.fields...)...)
}

func (w Wrapper) WithFields(err error, fields ...zap.Field) error {
	return WithFields(err, append(fields, w.fields...)...)
}

func (w Wrapper) WrapperWith(fields ...zap.Field) Wrapper {
	return Wrapper{fields: append(w.fields, fields...)}
}
