package timeutil

import (
	"fmt"
	"time"
)

const maxSec = int64(253402300799) // 9999-12-31T23:59:59Z

type Unit int

const (
	UnitSecond Unit = iota
	UnitMilli
	UnitMicro
	UnitNano
)

func Classify(unix int64) (Unit, error) {
	if unix < 0 {
		return UnitSecond, fmt.Errorf("timeutil: negative unix value %d not supported", unix)
	}
	if unix <= maxSec {
		return UnitSecond, nil
	}
	if unix <= maxSec*1e3 {
		return UnitMilli, nil
	}
	if unix <= maxSec*1e6 {
		return UnitMicro, nil
	}
	return UnitNano, nil
}

func ToSec(unix int64) (int64, error) {
	unit, err := Classify(unix)
	if err != nil {
		return 0, err
	}
	switch unit {
	case UnitSecond:
		return unix, nil
	case UnitMilli:
		return unix / 1e3, nil
	case UnitMicro:
		return unix / 1e6, nil
	case UnitNano:
		return unix / 1e9, nil
	default:
		return unix, nil
	}
}

func ToMilli(unix int64) (int64, error) {
	unit, err := Classify(unix)
	if err != nil {
		return 0, err
	}
	switch unit {
	case UnitSecond:
		return unix * 1e3, nil
	case UnitMilli:
		return unix, nil
	case UnitMicro:
		return unix / 1e3, nil
	case UnitNano:
		return unix / 1e6, nil
	default:
		return unix, nil
	}
}

func ToTime(unix int64) (time.Time, error) {
	unit, err := Classify(unix)
	if err != nil {
		return time.Time{}, err
	}
	switch unit {
	case UnitSecond:
		return time.Unix(unix, 0).UTC(), nil
	case UnitMilli:
		return time.UnixMilli(unix).UTC(), nil
	case UnitMicro:
		return time.UnixMicro(unix).UTC(), nil
	case UnitNano:
		return time.Unix(0, unix).UTC(), nil
	default:
		return time.Unix(unix, 0).UTC(), nil
	}
}
