package CH_conn

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func ConnectCH() (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			//Database: "nomad_metrics",
			Username: "nomad",
			Password: "nomad",
		},
	})
	if err != nil {
		return nil, err
	}
	//v, err := conn.ServerVersion()
	//fmt.Println(v)

	return conn, nil
}