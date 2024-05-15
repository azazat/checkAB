package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if os.Args[1] == "job_list" {
		db, err := sql.Open("sqlite3", "/db/config.db")
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

		err = db.Ping()
		if err != nil {
			log.Fatal("connection failed", err)
		}

		query := `
			SELECT task_id,task_name
			FROM task_table`

		var rows *sql.Rows
		rows, err = db.Query(query)
		if err != nil {
			log.Fatal(err)
		}

		var a struct {
			b string
			c string
		}

		t := make([]map[string]string, 0)
		for rows.Next() {
			err := rows.Scan(&a.b, &a.c)
			if err != nil {
				log.Fatal(err)
			}
			t = append(t, map[string]string{"{#JOB.NAME}": a.c})
		}

		json, err := json.Marshal(t)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", string(json))
	}
	if os.Args[1] == "task_list" {
		db, err := sql.Open("sqlite3", "/db/config.db")
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

		err = db.Ping()
		if err != nil {
			log.Fatal("connection failed", err)
		}

		query := `
			SELECT task_id,task_name
			FROM task_table`

		var rows *sql.Rows

		rows, err = db.Query(query)
		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		type a struct {
			b string
			c string
		}

		taskListArr := make(map[string]map[string]int64)

		for rows.Next() {
			var q a

			err := rows.Scan(&q.b, &q.c)
			if err != nil {
				log.Fatal(err)
			}

			jobName := q.c
			taskId := q.b

			db2, err := sql.Open("sqlite3", "/db/activity.db")
			if err != nil {
				log.Fatal(err)
			}

			defer db2.Close()

			err = db2.Ping()
			if err != nil {
				log.Fatal("connection failed", err)
			}

			query := `
					SELECT status,time_start,time_end,job_action 
					FROM result_table
					WHERE task_id=$1
					ORDER BY result_id DESC
					LIMIT 1`

			var struc struct {
				status         int
				time_start     int
				stop_timestamp int64
				job_action     int
			}

			err = db2.QueryRow(query, taskId).Scan(&struc.status, &struc.time_start, &struc.stop_timestamp, &struc.job_action)

			if err != nil {
				log.Fatal(err)
			}

			intervalDay := int64(math.Floor(float64(((time.Now().Unix() - int64(struc.stop_timestamp)) / 86400))))

			if struc.job_action == 268435456 {
				taskListArr[jobName] = map[string]int64{
					"stop_timestamp": struc.stop_timestamp,
					"interval_day":   intervalDay,
					"status":         0,
				}

				fmt.Printf("str: %v", taskListArr)
			}

			if struc.status == 1 {
				taskListArr[jobName] = map[string]int64{
					"stop_timestamp": struc.stop_timestamp,
					"interval_day":   intervalDay,
					"status":         1,
				}

			}

			if struc.status == 3 {
				taskListArr[jobName] = map[string]int64{
					"stop_timestamp": struc.stop_timestamp,
					"interval_day":   intervalDay,
					"status":         3,
				}

			}

			if struc.status == 4 {
				taskListArr[jobName] = map[string]int64{
					"stop_timestamp": struc.stop_timestamp,
					"interval_day":   intervalDay,
					"status":         4,
				}

			}

			if struc.status == 2 && struc.job_action != 268435456 {
				taskListArr[jobName] = map[string]int64{
					"stop_timestamp": struc.stop_timestamp,
					"interval_day":   intervalDay,
					"status":         2,
				}

			}

			if struc.status == 5 {
				taskListArr[jobName] = map[string]int64{
					"stop_timestamp": struc.stop_timestamp,
					"interval_day":   intervalDay,
					"status":         5,
				}

			}
		}
		json, err := json.Marshal(taskListArr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", string(json))

	}
}
