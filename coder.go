package sutando

import (
	"fmt"
	"reflect"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type coder decimal.Decimal

func (c coder) EncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != _TYPE_DECIMAL {
		return bsoncodec.ValueEncoderError{
			Name:     "decimalEncodeValue",
			Types:    []reflect.Type{_TYPE_DECIMAL},
			Received: val,
		}
	}

	dec := val.Interface().(decimal.Decimal)
	return vw.WriteString(dec.String())
}

func (c coder) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.IsValid() || !val.CanSet() || val.Type() != _TYPE_DECIMAL {
		return bsoncodec.ValueDecoderError{
			Name:     "decimalDecodeValue",
			Types:    []reflect.Type{_TYPE_DECIMAL},
			Received: val,
		}
	}

	var value decimal.Decimal
	switch vr.Type() {
	case bsontype.Decimal128:
		dec, err := vr.ReadDecimal128()
		if err != nil {
			return err
		}
		value, err = decimal.NewFromString(dec.String())
		if err != nil {
			return err
		}
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil {
			return err
		}
		value, err = decimal.NewFromString(str)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("received invalid BSON type to decode into decimal.Decimal: %s", vr.Type())
	}

	val.Set(reflect.ValueOf(value))
	return nil
}
