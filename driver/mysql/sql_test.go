package mysql

import (
	"context"
	"fmt"
	"testing"

	"github.com/thinkgos/carp/driver"
)

func Test_Tidb_SQL_Parse(t *testing.T) {
	sql :=
		"CREATE TABLE `announce` (" +
			"`id` bigint NOT NULL AUTO_INCREMENT DEFAULT '111' COMMENT '主键'," +
			"`title` varchar(255) DEFAULT NULL COMMENT '标题'" +
			// "`content` varchar(2048) NOT NULL DEFAULT '{}' COMMENT '内容'," +
			// "`value1` float unsigned NOT NULL DEFAULT '1' COMMENT '值1'," +
			// "`value2` float(10,1) unsigned NOT NULL DEFAULT '2' COMMENT '值2'," +
			// "`value2` double(16,2) NOT NULL DEFAULT '3' COMMENT '值3'," +
			// "`value3` enum('00','SH') NOT NULL DEFAULT '4' COMMENT '值4'," +
			// "`created_at` datetime NOT NULL COMMENT '发布时间'," +
			// "`updated_at` datetime NOT NULL," +
			// "`deleted_at` bigint NOT NULL DEFAULT '0'," +
			// "PRIMARY KEY (`id`) USING BTREE," +
			// "KEY `uk_value1_value2` (`value1`,`value2`) USING BTREE," +
			// "CONSTRAINT `posts_ibfk_1` FOREIGN KEY (`value1`) REFERENCES `user` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT" +
			")ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='公告-面向所有人的消息';"

	d := &SQLTidb{}
	value, err := d.InspectSchema(context.Background(), &driver.InspectOption{Data: sql})
	if err != nil {
		t.Fatal(err)
	}
	_ = value
}

func Test_SQL_Parse(t *testing.T) {
	sql :=
		"CREATE TABLE `announce` (" +
			"`id` bigint NOT NULL AUTO_INCREMENT," +
			// "`title` varchar(255) NOT NULL COMMENT '标题'," +
			// "`content` varchar(2048) NOT NULL COMMENT '内容'," +
			"`value1` float unsigned NOT NULL DEFAULT '1' COMMENT '值1'," +
			"`value2` float(10,1) unsigned NOT NULL DEFAULT '2' COMMENT '值2'," +
			"`value2` double(16,2) NOT NULL DEFAULT '3' COMMENT '值3'," +
			"`value3` enum('00','SH') NOT NULL DEFAULT '4' COMMENT '值4'," +
			// "`created_at` datetime NOT NULL COMMENT '发布时间'," +
			// "`updated_at` datetime NOT NULL," +
			// "`deleted_at` bigint NOT NULL DEFAULT '0'," +
			"PRIMARY KEY (`id`)," +
			"KEY `uk_value1_value2` (`value1`,`value2`) USING BTREE" +
			")ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='公告-面向所有人的消息';"

	d := &SQL{}
	value, err := d.InspectSchema(context.Background(), &driver.InspectOption{Data: sql})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(value.Entities[0])
}
