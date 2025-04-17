package mysql

import (
	"github.com/thinkgos/enst/driver"
)

func init() {
	driver.RegisterDriver(driver.Mysql, &MySQL{})
	driver.RegisterDriver(driver.FileMysql, &SQL{})
	driver.RegisterDriver(driver.FileMysqlTidb, &SQLTidb{})
}
