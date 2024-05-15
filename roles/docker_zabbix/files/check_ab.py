#!/bin/python
import sqlite3
from datetime import datetime
from math import floor
import sys
import json

error = int()
jobs_list_arr = []
now_timestamp = int(datetime.utcnow().timestamp())

def IntervalDay(d1, d2):

    seconds = d1 - d2
    return floor(seconds / 86400)

try:
    # dbt = sqlite3.connect('/db/config.db')
    # dbj = sqlite3.connect('/db/activity.db')
    dbt = sqlite3.connect('config.db')
    dbj = sqlite3.connect('activity.db')
    cursor = dbt.cursor()
    cursor2 = dbj.cursor()

    if sys.argv[1] == 'job_list':
        
        sqlite_select_query = """SELECT task_id,task_name from task_table"""
        cursor.execute(sqlite_select_query)
        records = cursor.fetchall()
  
        for tasks in records:
            job_name = tasks[1]
            jobs_list_arr.append({"{#JOB.NAME}": job_name})

        cursor.close()
        out = json.dumps(jobs_list_arr)
        print(out)
        
    if sys.argv[1] == 'task_list':
    
        sqlite_select_query = """SELECT task_id,task_name from task_table"""
        cursor.execute(sqlite_select_query)
        records = cursor.fetchall()
        task_list_arr = dict()       
        
        for tasks in records:

            task_id = tasks[0]
            sqlite_select_query2 = "select * from result_table where task_id=%d order by result_id desc limit 1" % task_id
            cursor2.execute(sqlite_select_query2)
            jobs = cursor2.fetchall()
            start=datetime.utcfromtimestamp(jobs[0][5]).strftime("%d.%m.%Y %H:%M")
            stop=datetime.utcfromtimestamp(jobs[0][6]).strftime("%d.%m.%Y %H:%M")
            stop_timestamp = jobs[0][6]
            interval_day = IntervalDay(now_timestamp, stop_timestamp)
            job_name = tasks[1]

            # Not Backed up yet
            if jobs[0][13] == 268435456:
                task_list_arr[str(job_name)] = {'stop_timestamp': stop_timestamp, 'interval_day': interval_day, 'status': 0}
                error = 0
            
            # Running
            if jobs[0][1] == 1:
                task_list_arr[str(job_name)] = {'stop_timestamp': stop_timestamp, 'interval_day': interval_day, 'status': 1}
                error = 1

            # Warnings
            if jobs[0][1] == 3:
                task_list_arr[str(job_name)] = {'stop_timestamp': stop_timestamp, 'interval_day': interval_day, 'status': 3}
                error = 3
            
            # Error
            if jobs[0][1] == 4:
                task_list_arr[str(job_name)] = {'stop_timestamp': stop_timestamp, 'interval_day': interval_day, 'status': 4}
                error = 4
            
            # Successful
            if jobs[0][1] == 2 and jobs[0][13] != 268435456:
                task_list_arr[str(job_name)] = {'stop_timestamp': stop_timestamp, 'interval_day': interval_day, 'status': 2}
                error=2

            # Canceled
            if jobs[0][1] == 5:
                task_list_arr[str(job_name)] = {'stop_timestamp': stop_timestamp, 'interval_day': interval_day, 'status': 5}
                error = 5            

        out = json.dumps(task_list_arr)
        print(out)
    
except sqlite3.Error as error:
    print("Ошибка при работе c SQLite", error)

# finally:
#     if dbt:
#         dbt.close()
#         print("Соединение c SQLite закрыто")

sys.exit(error)