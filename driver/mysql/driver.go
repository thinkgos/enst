package mysql

import (
	"github.com/thinkgos/carp/driver"
)

func init() {
	driver.RegisterDriver(driver.Mysql, &MySQL{})
	driver.RegisterDriver(driver.FileMysql, &SQL{})
	driver.RegisterDriver(driver.FileMysqlTidb, &SQLTidb{})
}
