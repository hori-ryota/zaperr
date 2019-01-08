package zaperr

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zaperr struct {
	err    error
	fields []zapcore.Field
}

type fieldser interface {
	Fields() []zapcore.Field
}

type fieldsAppender interface {
	WithFields(...zapcore.Field)
}

type wrapper interface {
	Wrap(message string)
}

type wrapfer interface {
	Wrapf(format string, args ...interface{})
}

func (e zaperr) Error() string {
	return e.err.Error()
}

func (e zaperr) Fields() []zapcore.Field {
	return append(e.fields, zap.Error(e.err))
}

func (e *zaperr) WithFields(fields ...zapcore.Field) {
	if e == nil {
		return
	}
	e.fields = append(e.fields, fields...)
}

// Deprecated: rename to WithFields
func (e *zaperr) AppendFields(fields ...zapcore.Field) {
	e.WithFields(fields...)
}

// for github.com/pkg/errors
func (e zaperr) Cause() error {
	return e.err
}

// for github.com/pkg/errors
func (e *zaperr) Wrap(message string) {
	e.err = errors.Wrap(e.err, message)
}

// for github.com/pkg/errors
func (e *zaperr) Wrapf(format string, args ...interface{}) {
	e.err = errors.Wrapf(e.err, format, args...)
}

// for github.com/pkg/errors
func Wrap(err error, message string) error {
	if e, ok := err.(wrapper); ok {
		e.Wrap(message)
		return err
	}
	return &zaperr{
		err: errors.Wrap(err, message),
	}
}

// for github.com/pkg/errors
func Wrapf(err error, format string, args ...interface{}) error {
	if e, ok := err.(wrapfer); ok {
		e.Wrapf(format, args...)
		return err
	}
	return &zaperr{
		err: errors.Wrapf(err, format, args...),
	}
}

func Fields(err error) []zapcore.Field {
	if err == nil {
		return nil
	}
	if e, ok := err.(fieldser); ok {
		return e.Fields()
	}
	return []zapcore.Field{zap.Error(err)}
}

func WithFields(err error, fields ...zapcore.Field) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(fieldsAppender); ok {
		e.WithFields(fields...)
		return err
	}

	return &zaperr{
		err:    err,
		fields: fields,
	}
}

// Deprecated: rename to WithFields
func AppendFields(err error, fields ...zapcore.Field) error {
	return WithFields(err, fields...)
}

func ToField(err error) zapcore.Field {
	if err == nil {
		return zap.Skip()
	}
	if e, ok := err.(zapcore.ObjectMarshaler); ok {
		return zap.Object("error with fields", e)
	}
	return zap.Error(err)
}

func ToNamedField(name string, err error) zapcore.Field {
	if err == nil {
		return zap.Skip()
	}
	if e, ok := err.(zapcore.ObjectMarshaler); ok {
		return zap.Object(name, e)
	}
	return zap.NamedError(name, err)
}

func (e zaperr) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, f := range e.Fields() {
		f.AddTo(enc)
	}
	return nil
}
