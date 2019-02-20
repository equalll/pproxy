package otto
import "github.com/equalll/mydebug"

func toStringPrimitive(value Value) Value {mydebug.INFO()
	return _toPrimitive(value, defaultValueHintString)
}

func toNumberPrimitive(value Value) Value {mydebug.INFO()
	return _toPrimitive(value, defaultValueHintNumber)
}

func toPrimitive(value Value) Value {mydebug.INFO()
	return _toPrimitive(value, defaultValueNoHint)
}

func _toPrimitive(value Value, hint _defaultValueHint) Value {mydebug.INFO()
	switch value.kind {
	case valueNull, valueUndefined, valueNumber, valueString, valueBoolean:
		return value
	case valueObject:
		return value._object().DefaultValue(hint)
	}
	panic(hereBeDragons(value.kind, value))
}
