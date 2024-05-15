#!/bin/php
<?php

// open DBs
$dbt=new SQLite3("/db/config.db");
$dbj=new SQLite3("/db/activity.db");

// configs
$jobs_list_arr = array();
//$task_list_arr = array();

$now_timestamp = time();

/**
 * Функция считает количество дней между двумя датами
 *
 * @param string $d1 первая дата
 * @param string $d2 вторая дата
 *
 * @return number количество дней
 */
function IntervalDay($d1, $d2)
{
    $seconds = abs($d1 - $d2);
    return floor($seconds / 86400);
}

if ($argv[1] == 'job_list')
{
    $task=$dbt->query("select task_id,task_name from task_table");
    while($tasks=$task->fetchArray())
    {

        $job_id =  $tasks['task_id'];
        $job_name = $tasks['task_name'];

        $jobs_list_arr[]["{#JOB.NAME}"]="$job_name";

    }
    $out = json_encode($jobs_list_arr);
    echo $out;
}

if ($argv[1] == 'task_list')
{
    $task=$dbt->query("select task_id,task_name from task_table");

    while($tasks=$task->fetchArray())
    {
        // get last job of each task
        $job=$dbj->query("select * from result_table where task_id=$tasks[task_id] order by result_id desc limit 1");
        $jobs=$job->fetchArray();

        $start=date("d.m.Y H:m", $jobs['time_start']);
        $stop=date("d.m.Y H:m", $jobs['time_end']);
        $stop_timestamp = $jobs['time_end'];
        $interval_day = IntervalDay($now_timestamp, $stop_timestamp);

        $job_name = $tasks['task_name'];

        // Not Backed up yet
        if ($jobs['job_action'] == 268435456)
        {
            $task_list_arr["$job_name"]=array("stop_timestamp"=>$stop_timestamp,"interval_day"=>$interval_day,"status"=>0);
            $error=0;
        }

        // Running
        if ($jobs['status'] == 1)
        {
            $task_list_arr["$job_name"]=array("stop_timestamp"=>$stop_timestamp,"interval_day"=>$interval_day,"status"=>1);
            $error=1;
        }

        // Warnings
        if ($jobs['status'] == 3)
        {
            $task_list_arr["$job_name"]=array("stop_timestamp"=>$stop_timestamp,"interval_day"=>$interval_day,"status"=>3);
            $error=3;
        }

        // Error
        if ($jobs['status'] == 4)
        {
            $task_list_arr["$job_name"]=array("stop_timestamp"=>$stop_timestamp,"interval_day"=>$interval_day,"status"=>4);
            $error=4;
        }
        // Successful
        if (($jobs['status'] == 2) and ($jobs['job_action'] != 268435456))
        {
            $task_list_arr["$job_name"]=array("stop_timestamp"=>$stop_timestamp,"interval_day"=>$interval_day,"status"=>2);
            $error=2;
        }
        // Canceled
        if ($jobs['status'] == 5)
        {

        $task_list_arr["$job_name"]=array("stop_timestamp"=>$stop_timestamp,"interval_day"=>$interval_day,"status"=>5);
        $error=5;
        }
    }

$out = json_encode($task_list_arr);

echo $out;
exit($error);
}
?>