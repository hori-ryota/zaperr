package zaperr

import (
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
	AppendFields(...zapcore.Field)
}

func (e zaperr) Error() string {
	return e.err.Error()
}

func (e zaperr) Fields() []zapcore.Field {
	return append(e.fields, zap.Error(e.err))
}

func (e *zaperr) AppendFields(fields ...zapcore.Field) {
	if e == nil {
		return
	}
	e.fields = append(e.fields, fields...)
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

func AppendFields(err error, fields ...zapcore.Field) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(fieldsAppender); ok {
		e.AppendFields(fields...)
		return err
	}

	return &zaperr{
		err:    err,
		fields: fields,
	}
}

func ToField(name string, err error) zapcore.Field {
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
