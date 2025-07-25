package logger

import "go.uber.org/zap"

func convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		switch v := f.Value.(type) {
		case string:
			zapFields[i] = zap.String(f.Key, v)
		case int:
			zapFields[i] = zap.Int(f.Key, v)
		case int64:
			zapFields[i] = zap.Int64(f.Key, v)
		case float64:
			zapFields[i] = zap.Float64(f.Key, v)
		case bool:
			zapFields[i] = zap.Bool(f.Key, v)
		case error:
			zapFields[i] = zap.Error(v)
		default:
			zapFields[i] = zap.Any(f.Key, v)
		}
	}
	return zapFields
}

//Поля для конвертации
var (
	String  = func(key, value string) Field { return Field{Key: key, Value: value} }
	Int     = func(key string, value int) Field { return Field{Key: key, Value: value} }
	Int64   = func(key string, value int64) Field { return Field{Key: key, Value: value} }
	Float64 = func(key string, value float64) Field { return Field{Key: key, Value: value} }
	Bool    = func(key string, value bool) Field { return Field{Key: key, Value: value} }
	Error   = func(err error) Field { return Field{Key: "error", Value: err} }
	Any     = func(key string, value interface{}) Field { return Field{Key: key, Value: value} }
)
