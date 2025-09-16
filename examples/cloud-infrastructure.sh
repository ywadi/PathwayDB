#!/bin/bash

# Cloud Infrastructure Sample Data
# Creates a comprehensive cloud architecture with multi-region deployment, security, and monitoring

echo "☁️ Creating cloud infrastructure graph..."

# Create the main graph
redis-cli -p 6379 GRAPH.CREATE cloud_infrastructure "Multi-cloud enterprise infrastructure"

echo "Creating cloud providers and regions..."
redis-cli -p 6379 NODE.CREATE cloud_infrastructure aws-us-east provider '{"name":"AWS US East","provider":"amazon_web_services","region":"us-east-1","availability_zones":["1a","1b","1c"],"compliance":["SOC2","PCI","HIPAA"],"cost_center":"primary"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure aws-eu-west provider '{"name":"AWS EU West","provider":"amazon_web_services","region":"eu-west-1","availability_zones":["1a","1b","1c"],"compliance":["GDPR","SOC2"],"cost_center":"emea"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure azure-central provider '{"name":"Azure Central US","provider":"microsoft_azure","region":"centralus","availability_zones":["1","2","3"],"compliance":["SOC2","ISO27001"],"cost_center":"hybrid"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure gcp-us-west provider '{"name":"GCP US West","provider":"google_cloud","region":"us-west1","availability_zones":["a","b","c"],"compliance":["SOC2","ISO27001"],"cost_center":"analytics"}'

echo "Creating compute resources..."
redis-cli -p 6379 NODE.CREATE cloud_infrastructure web-cluster compute '{"name":"Web Application Cluster","type":"kubernetes","nodes":12,"instance_type":"m5.xlarge","auto_scaling":"enabled","min_nodes":6,"max_nodes":24,"cpu_utilization_target":"70%"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure api-cluster compute '{"name":"API Services Cluster","type":"kubernetes","nodes":8,"instance_type":"c5.2xlarge","auto_scaling":"enabled","min_nodes":4,"max_nodes":16,"memory_utilization_target":"80%"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure worker-cluster compute '{"name":"Background Workers","type":"kubernetes","nodes":6,"instance_type":"m5.large","auto_scaling":"enabled","min_nodes":3,"max_nodes":12,"queue_length_target":"100"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure ml-cluster compute '{"name":"ML Training Cluster","type":"kubernetes","nodes":4,"instance_type":"p3.2xlarge","auto_scaling":"scheduled","gpu_enabled":true,"spot_instances":"80%","training_schedule":"nightly"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure edge-nodes compute '{"name":"Edge Computing Nodes","type":"ec2","count":50,"instance_type":"t3.micro","regions":["global"],"latency_target":"<50ms","content_caching":"enabled"}'

echo "Creating storage services..."
redis-cli -p 6379 NODE.CREATE cloud_infrastructure object-storage storage '{"name":"Object Storage","type":"s3","size":"500TB","storage_classes":["standard","ia","glacier"],"versioning":"enabled","encryption":"sse_s3","lifecycle_policies":"automated"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure block-storage storage '{"name":"Block Storage","type":"ebs","size":"100TB","volume_type":"gp3","iops":"10000","throughput":"1000MB/s","snapshots":"daily","encryption":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure file-storage storage '{"name":"Shared File Storage","type":"efs","size":"50TB","performance_mode":"general_purpose","throughput_mode":"provisioned","backup":"enabled","access_points":25}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure cache-storage storage '{"name":"Distributed Cache","type":"elasticache_redis","cluster_mode":"enabled","node_type":"r6g.xlarge","nodes":6,"memory":"96GB","backup":"enabled"}'

echo "Creating database services..."
redis-cli -p 6379 NODE.CREATE cloud_infrastructure primary-db database '{"name":"Primary Database","type":"rds_postgresql","instance_class":"db.r5.4xlarge","storage":"2TB","multi_az":"enabled","read_replicas":3,"backup_retention":"30_days","encryption":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure analytics-db database '{"name":"Analytics Database","type":"redshift","node_type":"dc2.8xlarge","nodes":8,"storage":"16TB","compression":"enabled","columnar_storage":"optimized","query_performance":"high"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure nosql-db database '{"name":"NoSQL Database","type":"dynamodb","capacity_mode":"on_demand","global_tables":"enabled","point_in_time_recovery":"enabled","encryption":"customer_managed","streams":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure time-series-db database '{"name":"Time Series Database","type":"timestream","memory_store_retention":"24h","magnetic_store_retention":"1_year","query_engine":"optimized","compression":"automatic"}'

echo "Creating networking components..."
redis-cli -p 6379 NODE.CREATE cloud_infrastructure vpc-primary network '{"name":"Primary VPC","type":"vpc","cidr":"10.0.0.0/16","subnets":{"public":3,"private":6,"database":3},"nat_gateways":3,"internet_gateway":"enabled","flow_logs":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure load-balancer network '{"name":"Application Load Balancer","type":"alb","scheme":"internet_facing","listeners":["HTTP","HTTPS"],"ssl_termination":"enabled","waf":"enabled","sticky_sessions":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure cdn network '{"name":"Content Delivery Network","type":"cloudfront","edge_locations":"global","cache_behaviors":10,"origin_shield":"enabled","compression":"enabled","http2":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure api-gateway network '{"name":"API Gateway","type":"api_gateway_v2","protocol":"HTTP","throttling":"10000_req/sec","caching":"enabled","authorization":"cognito","cors":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure vpn-gateway network '{"name":"VPN Gateway","type":"site_to_site_vpn","connections":5,"bandwidth":"1Gbps","redundancy":"enabled","routing":"bgp","encryption":"ipsec"}'

echo "Creating security services..."
redis-cli -p 6379 NODE.CREATE cloud_infrastructure waf security '{"name":"Web Application Firewall","type":"aws_waf","rules":50,"rate_limiting":"enabled","geo_blocking":"enabled","sql_injection_protection":"enabled","xss_protection":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure secrets-manager security '{"name":"Secrets Manager","type":"aws_secrets_manager","secrets":200,"rotation":"automatic","encryption":"kms","cross_region_replication":"enabled","audit_logging":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure identity-provider security '{"name":"Identity Provider","type":"cognito","user_pools":5,"identity_pools":3,"mfa":"enabled","social_providers":["google","facebook"],"saml":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure certificate-manager security '{"name":"Certificate Manager","type":"acm","certificates":25,"auto_renewal":"enabled","domain_validation":"dns","wildcard_support":"enabled","transparency_logging":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure security-hub security '{"name":"Security Hub","type":"aws_security_hub","standards":["aws_foundational","cis","pci_dss"],"findings_aggregation":"enabled","custom_insights":"enabled","compliance_score":"tracked"}'

echo "Creating monitoring and logging..."
redis-cli -p 6379 NODE.CREATE cloud_infrastructure cloudwatch monitoring '{"name":"CloudWatch","type":"aws_cloudwatch","metrics":"custom_enabled","logs_retention":"1_year","alarms":500,"dashboards":50,"insights_queries":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure xray monitoring '{"name":"X-Ray Tracing","type":"aws_xray","sampling_rate":"10%","trace_retention":"30_days","service_map":"enabled","annotations":"custom","insights":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure elk-stack monitoring '{"name":"ELK Stack","type":"elasticsearch_service","cluster_size":6,"storage":"10TB","log_retention":"90_days","kibana_users":100,"alerting":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure prometheus monitoring '{"name":"Prometheus","type":"managed_prometheus","retention":"1_year","ingestion_rate":"1M_samples/sec","query_performance":"optimized","alertmanager":"integrated"}'

echo "Creating DevOps and CI/CD..."
redis-cli -p 6379 NODE.CREATE cloud_infrastructure ci-cd devops '{"name":"CI/CD Pipeline","type":"codepipeline","pipelines":25,"stages_per_pipeline":"avg_5","deployment_frequency":"daily","success_rate":"98%","rollback_capability":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure container-registry devops '{"name":"Container Registry","type":"ecr","repositories":100,"image_scanning":"enabled","lifecycle_policies":"enabled","cross_region_replication":"enabled","vulnerability_scanning":"continuous"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure infrastructure-as-code devops '{"name":"Infrastructure as Code","type":"terraform_cloud","workspaces":50,"state_management":"remote","policy_enforcement":"enabled","cost_estimation":"enabled","drift_detection":"enabled"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure config-management devops '{"name":"Configuration Management","type":"aws_config","rules":100,"compliance_tracking":"enabled","remediation":"automatic","change_tracking":"comprehensive","resource_inventory":"real_time"}'

echo "Creating backup and disaster recovery..."
redis-cli -p 6379 NODE.CREATE cloud_infrastructure backup-service backup '{"name":"Backup Service","type":"aws_backup","backup_plans":20,"retention_policies":"tiered","cross_region_backup":"enabled","compliance":"enabled","restore_testing":"automated"}'
redis-cli -p 6379 NODE.CREATE cloud_infrastructure disaster-recovery backup '{"name":"Disaster Recovery","type":"multi_region","rto":"4_hours","rpo":"15_minutes","failover":"automated","testing_frequency":"quarterly","documentation":"maintained"}'

echo "Creating compute-to-network connections..."
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure web-to-vpc web-cluster vpc-primary deploys_in '{"subnet_type":"private","security_groups":["web_sg","common_sg"],"network_acls":"restrictive","outbound_rules":"https_only"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure api-to-vpc api-cluster vpc-primary deploys_in '{"subnet_type":"private","security_groups":["api_sg","common_sg"],"network_acls":"restrictive","outbound_rules":"database_cache"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure lb-to-web load-balancer web-cluster routes_to '{"target_group":"web_servers","health_check":"enabled","stickiness":"enabled","algorithm":"round_robin","ssl_termination":"lb_level"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure cdn-to-lb cdn load-balancer origins_from '{"cache_behavior":"default","ttl":"1_hour","compression":"enabled","viewer_protocol_policy":"redirect_to_https","origin_shield":"enabled"}'

echo "Creating storage connections..."
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure web-to-object web-cluster object-storage uses '{"access_pattern":"read_write","permissions":"iam_role","encryption":"transit_rest","lifecycle":"intelligent_tiering","versioning":"enabled"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure api-to-cache api-cluster cache-storage connects_to '{"connection_pooling":"enabled","ssl":"enabled","cluster_mode":"enabled","failover":"automatic","monitoring":"cloudwatch"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure worker-to-file worker-cluster file-storage mounts '{"mount_point":"/shared","performance_mode":"general_purpose","access_point":"workers","encryption":"transit_rest"}'

echo "Creating database connections..."
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure api-to-primary api-cluster primary-db queries '{"connection_pooling":"pgbouncer","ssl":"required","read_write_split":"enabled","monitoring":"performance_insights","backup":"point_in_time"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure analytics-to-redshift ml-cluster analytics-db analyzes '{"connection_type":"jdbc","ssl":"enabled","query_optimization":"enabled","workload_management":"enabled","result_caching":"enabled"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure api-to-nosql api-cluster nosql-db accesses '{"access_pattern":"on_demand","consistency":"eventual","global_tables":"enabled","streams":"kinesis","encryption":"customer_managed"}'

echo "Creating security connections..."
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure lb-to-waf load-balancer waf protects_with '{"rule_evaluation":"priority_based","logging":"enabled","metrics":"cloudwatch","geo_restriction":"enabled","rate_limiting":"per_ip"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure api-to-secrets api-cluster secrets-manager retrieves_from '{"rotation":"automatic","encryption":"kms","cross_region":"enabled","audit_trail":"cloudtrail","access_logging":"enabled"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure web-to-cognito web-cluster identity-provider authenticates_with '{"user_pools":"enabled","identity_pools":"enabled","mfa":"required","social_providers":"enabled","saml":"enterprise"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure lb-to-acm load-balancer certificate-manager uses_certificates_from '{"auto_renewal":"enabled","domain_validation":"dns","wildcard":"enabled","transparency_logging":"enabled"}'

echo "Creating monitoring connections..."
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure all-to-cloudwatch web-cluster cloudwatch monitors_with '{"metrics":"custom","logs":"application","alarms":"threshold_based","dashboards":"operational","insights":"enabled"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure api-to-xray api-cluster xray traces_with '{"sampling_rate":"10%","annotations":"custom","subsegments":"detailed","service_map":"enabled","insights":"performance"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure logs-to-elk cloudwatch elk-stack forwards_to '{"log_groups":"application","retention":"90_days","parsing":"json","indexing":"time_based","alerting":"threshold"}'

echo "Creating DevOps connections..."
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure cicd-to-registry ci-cd container-registry builds_to '{"image_scanning":"enabled","vulnerability_assessment":"high_critical","lifecycle_policies":"enabled","cross_region_replication":"enabled"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure registry-to-cluster container-registry web-cluster deploys_to '{"image_pull_policy":"always","security_scanning":"required","resource_limits":"enforced","rolling_updates":"enabled"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure iac-to-config infrastructure-as-code config-management manages_with '{"drift_detection":"enabled","compliance_rules":"enforced","change_tracking":"comprehensive","remediation":"automatic"}'

echo "Creating backup connections..."
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure db-to-backup primary-db backup-service backs_up_to '{"frequency":"daily","retention":"30_days","cross_region":"enabled","encryption":"enabled","restore_testing":"monthly"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure backup-to-dr backup-service disaster-recovery replicates_to '{"rto":"4_hours","rpo":"15_minutes","automated_failover":"enabled","testing":"quarterly","documentation":"current"}'

echo "Creating multi-region connections..."
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure primary-to-secondary aws-us-east aws-eu-west replicates_to '{"replication_type":"async","data_sync":"continuous","failover_capability":"manual","compliance":"gdpr","latency":"<100ms"}'
redis-cli -p 6379 EDGE.CREATE cloud_infrastructure hybrid-connection aws-us-east azure-central connects_via '{"connection_type":"express_route","bandwidth":"10Gbps","redundancy":"enabled","encryption":"ipsec","routing":"bgp"}'

echo "✅ Cloud infrastructure graph created successfully!"
echo "📊 Graph contains:"
echo "   - 4 Cloud provider regions with compliance frameworks"
echo "   - 5 Compute clusters with auto-scaling capabilities"
echo "   - 4 Storage services with different access patterns"
echo "   - 4 Database services optimized for different workloads"
echo "   - 5 Networking components for traffic management"
echo "   - 5 Security services for comprehensive protection"
echo "   - 4 Monitoring and observability tools"
echo "   - 4 DevOps and automation services"
echo "   - 2 Backup and disaster recovery solutions"
echo "   - 25+ edges showing infrastructure relationships"
echo ""
echo "🔄 Infrastructure Flow Examples Created:"
echo "   1. Web Traffic: CDN → Load Balancer → WAF → Web Cluster → VPC"
echo "   2. Data Flow: API Cluster → Primary DB → Backup → Disaster Recovery"
echo "   3. Security: Identity Provider → Secrets Manager → Certificate Manager → Security Hub"
echo "   4. DevOps: CI/CD → Container Registry → Kubernetes → Monitoring"
echo "   5. Multi-Cloud: AWS ↔ Azure ↔ GCP (hybrid connectivity)"
