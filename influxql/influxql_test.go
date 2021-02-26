// Copyright 2016 Eleme. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package influxql

import "testing"

func TestInfluxQL(t *testing.T) {
	checkPoint(t, "select * from cpu", "cpu")
	checkPoint(t, "(select *) from cpu", "cpu")
	checkPoint(t, "[select *] from cpu", "cpu")
	checkPoint(t, "{select *} from cpu", "cpu")
	checkPoint(t, "select * from \"cpu\"", "cpu")
	checkPoint(t, "select * from \"c\\\"pu\"", "c\"pu")
	checkPoint(t, "select * from 'cpu'", "cpu")

	checkPoint(t, "SELECT mean(\"value\") FROM \"cpu\" WHERE \"region\" = 'uswest' GROUP BY time(10m) fill(0)", "cpu")
	checkPoint(t, "SELECT mean(\"value\") INTO \"cpu\\\"_1h\".:MEASUREMENT FROM /cpu.*/", "/cpu.*/")

	checkPoint(t, "REVOKE ALL PRIVILEGES FROM \"jdoe\"", "jdoe")
	checkPoint(t, "REVOKE READ ON \"mydb\" FROM \"jdoe\"", "jdoe")

	checkPoint(t, "DELETE FROM \"cpu\"", "cpu")
	checkPoint(t, "DELETE FROM \"cpu\" WHERE time < '2000-01-01T00:00:00Z'", "cpu")

	// checkPoint(t, "DROP SERIES FROM \"telegraf\".\"autogen\".\"cpu\" WHERE cpu = 'cpu8'", "cpu")
	// checkPoint(t, "SHOW FIELD KEYS", "cpu")
	checkPoint(t, "SHOW FIELD KEYS FROM \"cpu\"", "cpu")
	// checkPoint(t, "SHOW SERIES FROM \"telegraf\".\"autogen\".\"cpu\" WHERE cpu = 'cpu8'", "cpu")

	// checkPoint(t, "SHOW TAG KEYS", "cpu")
	checkPoint(t, "SHOW TAG KEYS FROM cpu", "cpu")
	checkPoint(t, "SHOW TAG KEYS FROM \"cpu\" WHERE \"region\" = 'uswest'", "cpu")
	// checkPoint(t, "SHOW TAG KEYS WHERE \"host\" = 'serverA'", "cpu")

	// checkPoint(t, "SHOW TAG VALUES WITH KEY = \"region\"", "cpu")
	checkPoint(t, "SHOW TAG VALUES FROM \"cpu\" WITH KEY = \"region\"", "cpu")
	// checkPoint(t, "SHOW TAG VALUES WITH KEY !~ /.*c.*/", "cpu")
	checkPoint(t, "SHOW TAG VALUES FROM \"cpu\" WITH KEY IN (\"region\", \"host\") WHERE \"service\" = 'redis'", "cpu")

	checkPoint(t, "SHOW FIELD KEYS FROM \"1h\".\"cpu\"", "cpu")
	checkPoint(t, "SHOW FIELD KEYS FROM 1h.cpu", "cpu")
	checkPoint(t, "SHOW FIELD KEYS FROM \"cpu.load\"", "cpu.load")
	checkPoint(t, "SHOW FIELD KEYS FROM 1h.\"cpu.load\"", "cpu.load")
	checkPoint(t, "SHOW FIELD KEYS FROM \"1h\".\"cpu.load\"", "cpu.load")
}

func checkPoint(t *testing.T, q string, m string) {
	qm, err := GetMeasurementFromInfluxQL(q)
	if err != nil {
		t.Errorf("error: %s", err)
		return
	}
	if qm != m {
		t.Errorf("measurement wrong: %s != %s", qm, m)
		return
	}
}

func BenchmarkInfluxQL(b *testing.B) {
	q := "SELECT mean(\"value\") FROM \"cpu\" WHERE \"region\" = 'uswest' GROUP BY time(10m) fill(0)"
	for i := 0; i < b.N; i++ {
		qm, err := GetMeasurementFromInfluxQL(q)
		if err != nil {
			b.Errorf("error: %s", err)
			return
		}
		if qm != "cpu" {
			b.Errorf("measurement wrong: %s != %s", qm, "cpu")
			return
		}
	}
}
