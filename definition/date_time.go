package definition

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"haibara/consts"
	"reflect"
	"time"
)

type DateTime time.Time

func (dateTime *DateTime) Scan(v any) error {
	newTime, ok := v.(time.Time)
	if !ok {
		return errors.New(fmt.Sprintf("不能将%v转换为DateTime", reflect.TypeOf(v)))
	}
	*dateTime = DateTime(newTime)
	return nil
}

func (dateTime *DateTime) Value() (driver.Value, error) {
	return time.Time(*dateTime), nil
}

func (dateTime *DateTime) MarshalJSON() ([]byte, error) {
	timeString := fmt.Sprintf(`"%s"`, time.Time(*dateTime).Format(consts.DateTime))
	return []byte(timeString), nil
}

func (dateTime *DateTime) UnmarshalJSON(data []byte) (err error) {
	newTime, err := time.Parse(consts.DateTime, string(data))
	if err != nil {
		return err
	}
	*dateTime = DateTime(newTime)
	return
}

func (dateTime *DateTime) String() string {
	return time.Time(*dateTime).Format(consts.DateTime)
}
