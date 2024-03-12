/*  Clickhouse nomad.metrics schema file    */

CREATE DATABASE IF NOT EXISTS nomad;
CREATE table IF NOT EXISTS nomad.metrics(
    Source      String ,
    MetricType String  COMMENT 'counter or gauge',
	StatName   String ,
	MetricName String ,
    Timestamp DateTime64(3, 'Asia/Calcutta'),
	StatValue  Float64,
    Labels Map(String,String)
)ENGINE = MergeTree()
ORDER BY (Source,MetricType,StatName,MetricName,Timestamp,StatValue)
PARTITION BY toYYYYMM(Timestamp)
PRIMARY KEY (Source,MetricType,StatName,MetricName,Timestamp);
/*
    TOdO: Skipping indexes : only applicable on MergeTree based tables
    - these are called skipping because they enable clickhouse to skip reading 
    significant amounts of data that are "guaranteed" to have no matching values
    -
*/