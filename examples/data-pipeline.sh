#!/bin/bash

# Data Pipeline and Analytics Sample Data
# Creates a comprehensive data processing pipeline with sources, transformations, and destinations

echo "📊 Creating data pipeline and analytics graph..."

# Create the main graph
redis-cli -p 6379 GRAPH.CREATE data_pipeline "Enterprise data pipeline and analytics platform"

echo "Creating data sources..."
redis-cli -p 6379 NODE.CREATE data_pipeline customer-db source '{"name":"Customer Database","type":"postgresql","host":"prod-db-01","port":5432,"schema":"customers","tables":["users","profiles","preferences"],"size":"500GB","update_frequency":"real_time"}'
redis-cli -p 6379 NODE.CREATE data_pipeline sales-api source '{"name":"Sales API","type":"rest_api","endpoint":"https://api.sales.com/v2","rate_limit":"1000/min","data_format":"json","authentication":"oauth2","refresh_interval":"5m"}'
redis-cli -p 6379 NODE.CREATE data_pipeline web-analytics source '{"name":"Web Analytics","type":"google_analytics","property_id":"GA4-12345","metrics":["sessions","pageviews","conversions"],"dimensions":["source","medium","campaign"],"sampling":"unsampled"}'
redis-cli -p 6379 NODE.CREATE data_pipeline mobile-events source '{"name":"Mobile App Events","type":"event_stream","platform":"firebase","events":["app_open","purchase","user_engagement"],"volume":"1M_events/day","retention":"90_days"}'
redis-cli -p 6379 NODE.CREATE data_pipeline social-media source '{"name":"Social Media Data","type":"social_api","platforms":["twitter","facebook","instagram"],"metrics":["engagement","reach","mentions"],"sentiment_analysis":true,"update_frequency":"hourly"}'
redis-cli -p 6379 NODE.CREATE data_pipeline iot-sensors source '{"name":"IoT Sensor Data","type":"mqtt_stream","devices":500,"metrics":["temperature","humidity","pressure"],"frequency":"1_reading/second","protocol":"mqtt_v3.1.1"}'

echo "Creating data ingestion layer..."
redis-cli -p 6379 NODE.CREATE data_pipeline kafka-ingestion ingestion '{"name":"Kafka Data Ingestion","type":"apache_kafka","cluster_size":3,"topics":12,"partitions_per_topic":6,"replication_factor":3,"throughput":"100K_msg/sec"}'
redis-cli -p 6379 NODE.CREATE data_pipeline api-gateway ingestion '{"name":"API Gateway","type":"kong","rate_limiting":"10K_req/min","authentication":["api_key","oauth2","jwt"],"data_validation":true,"request_logging":true}'
redis-cli -p 6379 NODE.CREATE data_pipeline file-ingestion ingestion '{"name":"File Ingestion Service","type":"custom_service","supported_formats":["csv","json","parquet","avro"],"max_file_size":"10GB","parallel_processing":true}'
redis-cli -p 6379 NODE.CREATE data_pipeline stream-processor ingestion '{"name":"Stream Processor","type":"apache_flink","parallelism":16,"checkpointing":"enabled","state_backend":"rocksdb","latency":"sub_second"}'

echo "Creating data transformation layer..."
redis-cli -p 6379 NODE.CREATE data_pipeline etl-pipeline transformation '{"name":"ETL Pipeline","type":"apache_airflow","dags":25,"tasks_per_dag":"avg_8","schedule":"various","retry_policy":"3_attempts","monitoring":"enabled"}'
redis-cli -p 6379 NODE.CREATE data_pipeline data-quality transformation '{"name":"Data Quality Engine","type":"great_expectations","validation_rules":150,"data_profiling":true,"anomaly_detection":true,"alert_threshold":"5%_failure_rate"}'
redis-cli -p 6379 NODE.CREATE data_pipeline feature-store transformation '{"name":"Feature Store","type":"feast","features":500,"feature_groups":50,"serving_latency":"<10ms","batch_scoring":"hourly","real_time_serving":true}'
redis-cli -p 6379 NODE.CREATE data_pipeline ml-pipeline transformation '{"name":"ML Pipeline","type":"kubeflow","models":15,"training_frequency":"weekly","model_validation":"automated","a_b_testing":"enabled","model_registry":"mlflow"}'

echo "Creating data storage layer..."
redis-cli -p 6379 NODE.CREATE data_pipeline data-lake storage '{"name":"Data Lake","type":"amazon_s3","size":"50TB","storage_classes":["standard","ia","glacier"],"lifecycle_policies":"enabled","encryption":"sse_s3","versioning":"enabled"}'
redis-cli -p 6379 NODE.CREATE data_pipeline data-warehouse storage '{"name":"Data Warehouse","type":"snowflake","size":"10TB","compute_clusters":5,"auto_scaling":"enabled","query_acceleration":"enabled","time_travel":"90_days"}'
redis-cli -p 6379 NODE.CREATE data_pipeline operational-db storage '{"name":"Operational Database","type":"mongodb","cluster_size":3,"sharding":"enabled","replica_set":"rs0","storage_engine":"wiredTiger","backup_frequency":"daily"}'
redis-cli -p 6379 NODE.CREATE data_pipeline cache-layer storage '{"name":"Cache Layer","type":"redis_cluster","nodes":6,"memory":"64GB","persistence":"rdb","eviction_policy":"allkeys_lru","high_availability":"enabled"}'
redis-cli -p 6379 NODE.CREATE data_pipeline search-index storage '{"name":"Search Index","type":"elasticsearch","cluster_size":5,"indices":20,"shards_per_index":3,"replicas":1,"query_performance":"<100ms"}'

echo "Creating analytics and BI layer..."
redis-cli -p 6379 NODE.CREATE data_pipeline bi-platform analytics '{"name":"Business Intelligence Platform","type":"tableau","users":200,"dashboards":150,"data_sources":25,"refresh_frequency":"hourly","mobile_support":true}'
redis-cli -p 6379 NODE.CREATE data_pipeline ad-hoc-analytics analytics '{"name":"Ad-hoc Analytics","type":"jupyter_hub","users":50,"notebooks":500,"kernels":["python","r","scala"],"compute_resources":"kubernetes","collaboration":"enabled"}'
redis-cli -p 6379 NODE.CREATE data_pipeline real-time-dashboard analytics '{"name":"Real-time Dashboard","type":"grafana","dashboards":75,"data_sources":15,"alert_rules":200,"update_frequency":"1_second","user_access":"role_based"}'
redis-cli -p 6379 NODE.CREATE data_pipeline reporting-engine analytics '{"name":"Reporting Engine","type":"jasperreports","reports":300,"schedules":150,"output_formats":["pdf","excel","csv"],"distribution":"email_automated"}'

echo "Creating data governance layer..."
redis-cli -p 6379 NODE.CREATE data_pipeline data-catalog governance '{"name":"Data Catalog","type":"apache_atlas","assets":5000,"lineage_tracking":"automated","metadata_management":"comprehensive","search_capabilities":"full_text"}'
redis-cli -p 6379 NODE.CREATE data_pipeline access-control governance '{"name":"Access Control","type":"apache_ranger","policies":500,"user_groups":50,"audit_logging":"comprehensive","fine_grained_permissions":true}'
redis-cli -p 6379 NODE.CREATE data_pipeline privacy-compliance governance '{"name":"Privacy Compliance","type":"privacera","gdpr_compliance":true,"data_masking":"dynamic","consent_management":"automated","breach_detection":"enabled"}'
redis-cli -p 6379 NODE.CREATE data_pipeline data-lineage governance '{"name":"Data Lineage","type":"datahub","entities":10000,"relationships":50000,"impact_analysis":"automated","change_tracking":"real_time"}'

echo "Creating monitoring and observability..."
redis-cli -p 6379 NODE.CREATE data_pipeline pipeline-monitoring monitoring '{"name":"Pipeline Monitoring","type":"datadog","metrics":1000,"alerts":200,"dashboards":50,"sla_tracking":"enabled","anomaly_detection":"ml_based"}'
redis-cli -p 6379 NODE.CREATE data_pipeline cost-optimization monitoring '{"name":"Cost Optimization","type":"kubecost","cost_allocation":"granular","budget_alerts":"enabled","rightsizing_recommendations":"automated","waste_detection":"continuous"}'
redis-cli -p 6379 NODE.CREATE data_pipeline performance-monitoring monitoring '{"name":"Performance Monitoring","type":"new_relic","apm":"enabled","infrastructure_monitoring":"comprehensive","log_analysis":"centralized","trace_analysis":"distributed"}'

echo "Creating data flow connections - Ingestion Layer..."
redis-cli -p 6379 EDGE.CREATE data_pipeline customer-to-kafka customer-db kafka-ingestion streams '{"connector":"debezium","format":"avro","schema_registry":"confluent","change_capture":"real_time","compression":"snappy"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline api-to-gateway sales-api api-gateway ingests '{"rate_limiting":"applied","authentication":"validated","payload_transformation":"json_to_avro","error_handling":"retry_with_backoff"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline mobile-to-stream mobile-events stream-processor processes '{"windowing":"tumbling_5min","aggregations":["count","sum","avg"],"late_data_handling":"allowed_1min","state_management":"keyed"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline iot-to-kafka iot-sensors kafka-ingestion streams '{"topic":"iot_readings","partitioning":"device_id","serialization":"protobuf","batch_size":"1000_messages","compression":"lz4"}'

echo "Creating transformation connections..."
redis-cli -p 6379 EDGE.CREATE data_pipeline kafka-to-etl kafka-ingestion etl-pipeline consumes '{"consumer_group":"etl_processors","offset_management":"automatic","parallelism":"topic_partitions","error_handling":"dead_letter_queue"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline etl-to-quality etl-pipeline data-quality validates '{"validation_frequency":"every_batch","quality_metrics":["completeness","accuracy","consistency"],"threshold":"95%_pass_rate","remediation":"automated"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline quality-to-feature data-quality feature-store feeds '{"feature_validation":"enabled","feature_monitoring":"drift_detection","serving_preparation":"optimized","versioning":"semantic"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline feature-to-ml feature-store ml-pipeline trains '{"feature_selection":"automated","model_training":"scheduled","hyperparameter_tuning":"bayesian","cross_validation":"time_series_split"}'

echo "Creating storage connections..."
redis-cli -p 6379 EDGE.CREATE data_pipeline etl-to-lake etl-pipeline data-lake stores '{"format":"parquet","partitioning":"date_hour","compression":"snappy","metadata":"hive_metastore","lifecycle":"intelligent_tiering"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline lake-to-warehouse data-lake data-warehouse loads '{"schedule":"hourly","incremental_loading":"merge","data_validation":"row_count_checksum","performance":"query_optimization"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline warehouse-to-cache data-warehouse cache-layer caches '{"cache_strategy":"write_through","ttl":"1_hour","eviction":"lru","hit_ratio_target":"85%","warming_strategy":"predictive"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline processed-to-search etl-pipeline search-index indexes '{"indexing_strategy":"near_real_time","mapping":"dynamic","analyzers":"custom","query_optimization":"enabled"}'

echo "Creating analytics connections..."
redis-cli -p 6379 EDGE.CREATE data_pipeline warehouse-to-bi data-warehouse bi-platform queries '{"connection_pooling":"enabled","query_caching":"aggressive","performance_optimization":"materialized_views","security":"row_level"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline lake-to-jupyter data-lake ad-hoc-analytics analyzes '{"access_method":"s3_api","compute_engine":"spark","notebook_sharing":"collaborative","version_control":"git_integration"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline stream-to-grafana stream-processor real-time-dashboard visualizes '{"metrics_export":"prometheus","dashboard_updates":"real_time","alerting":"threshold_based","annotation":"event_driven"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline warehouse-to-reports data-warehouse reporting-engine generates '{"report_scheduling":"cron_based","parameter_passing":"dynamic","output_distribution":"multi_channel","template_management":"versioned"}'

echo "Creating governance connections..."
redis-cli -p 6379 EDGE.CREATE data_pipeline all-to-catalog data-warehouse data-catalog catalogs '{"metadata_harvesting":"automated","lineage_capture":"comprehensive","business_glossary":"maintained","data_discovery":"search_enabled"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline catalog-to-access data-catalog access-control enforces '{"policy_sync":"real_time","permission_inheritance":"hierarchical","audit_trail":"comprehensive","compliance_reporting":"automated"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline access-to-privacy access-control privacy-compliance integrates '{"data_classification":"automated","masking_policies":"dynamic","consent_enforcement":"real_time","breach_monitoring":"continuous"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline pipeline-to-lineage etl-pipeline data-lineage tracks '{"lineage_capture":"code_analysis","impact_analysis":"downstream_dependencies","change_propagation":"automated","visualization":"interactive"}'

echo "Creating monitoring connections..."
redis-cli -p 6379 EDGE.CREATE data_pipeline all-to-monitoring kafka-ingestion pipeline-monitoring observes '{"metrics_collection":"comprehensive","alert_configuration":"threshold_anomaly","dashboard_creation":"automated","sla_monitoring":"business_critical"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline monitoring-to-cost pipeline-monitoring cost-optimization optimizes '{"cost_attribution":"resource_tagging","budget_tracking":"real_time","optimization_recommendations":"ml_driven","waste_identification":"automated"}'
redis-cli -p 6379 EDGE.CREATE data_pipeline performance-monitors performance-monitoring stream-processor monitors '{"latency_tracking":"end_to_end","throughput_monitoring":"real_time","resource_utilization":"detailed","bottleneck_identification":"automated"}'

echo "✅ Data pipeline and analytics graph created successfully!"
echo "📊 Graph contains:"
echo "   - 6 Diverse data sources (DB, API, Analytics, Mobile, Social, IoT)"
echo "   - 4 Data ingestion services with different protocols"
echo "   - 4 Transformation services including ML and quality checks"
echo "   - 5 Storage layers from data lake to search index"
echo "   - 4 Analytics and BI platforms for different use cases"
echo "   - 4 Data governance and compliance tools"
echo "   - 3 Monitoring and observability solutions"
echo "   - 20+ edges showing complex data flow relationships"
echo ""
echo "🔄 Data Flow Examples Created:"
echo "   1. Batch Processing: DB → Kafka → ETL → Quality → Data Lake → Warehouse → BI"
echo "   2. Real-time Stream: IoT → Stream Processor → Feature Store → ML → Real-time Dashboard"
echo "   3. Governance Flow: All Data → Catalog → Access Control → Privacy Compliance"
echo "   4. Analytics Flow: Warehouse → BI/Jupyter/Reports for different user personas"
