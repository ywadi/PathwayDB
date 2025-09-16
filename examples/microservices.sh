#!/bin/bash

# Microservices Architecture Sample Data
# Creates a realistic microservices ecosystem with services, databases, and dependencies

echo "🏗️  Creating microservices architecture graph..."

# Create the main graph
redis-cli -p 6379 GRAPH.CREATE microservices "E-commerce microservices architecture"

echo "Creating API Gateway..."
redis-cli -p 6379 NODE.CREATE microservices api-gateway gateway '{"name":"API Gateway","type":"Kong","version":"3.0","port":8080,"load_balancer":"nginx"}'

echo "Creating core services..."
redis-cli -p 6379 NODE.CREATE microservices user-service service '{"name":"User Service","port":8001,"language":"Go","replicas":3,"cpu":"500m","memory":"1Gi","health_check":"/health"}'
redis-cli -p 6379 NODE.CREATE microservices auth-service service '{"name":"Authentication Service","port":8002,"language":"Go","replicas":2,"cpu":"300m","memory":"512Mi","jwt_secret":"encrypted"}'
redis-cli -p 6379 NODE.CREATE microservices product-service service '{"name":"Product Service","port":8003,"language":"Python","replicas":3,"cpu":"800m","memory":"2Gi","framework":"FastAPI"}'
redis-cli -p 6379 NODE.CREATE microservices order-service service '{"name":"Order Service","port":8004,"language":"Java","replicas":2,"cpu":"1000m","memory":"2Gi","framework":"Spring Boot"}'
redis-cli -p 6379 NODE.CREATE microservices payment-service service '{"name":"Payment Service","port":8005,"language":"Node.js","replicas":2,"cpu":"600m","memory":"1Gi","pci_compliant":true}'
redis-cli -p 6379 NODE.CREATE microservices notification-service service '{"name":"Notification Service","port":8006,"language":"Python","replicas":1,"cpu":"200m","memory":"256Mi","channels":["email","sms","push"]}'
redis-cli -p 6379 NODE.CREATE microservices inventory-service service '{"name":"Inventory Service","port":8007,"language":"Go","replicas":2,"cpu":"400m","memory":"1Gi","real_time":true}'

echo "Creating databases..."
redis-cli -p 6379 NODE.CREATE microservices user-db database '{"name":"User Database","type":"PostgreSQL","version":"15.0","size":"100GB","backup_schedule":"daily","encryption":true}'
redis-cli -p 6379 NODE.CREATE microservices product-db database '{"name":"Product Database","type":"MongoDB","version":"6.0","size":"500GB","sharding":true,"replica_set":"rs0"}'
redis-cli -p 6379 NODE.CREATE microservices order-db database '{"name":"Order Database","type":"PostgreSQL","version":"15.0","size":"200GB","partitioning":"monthly","wal_level":"replica"}'
redis-cli -p 6379 NODE.CREATE microservices payment-db database '{"name":"Payment Database","type":"PostgreSQL","version":"15.0","size":"50GB","encryption":"AES-256","audit_log":true}'
redis-cli -p 6379 NODE.CREATE microservices analytics-db database '{"name":"Analytics Database","type":"ClickHouse","version":"23.0","size":"1TB","compression":"zstd","retention":"2_years"}'

echo "Creating cache layers..."
redis-cli -p 6379 NODE.CREATE microservices redis-cache cache '{"name":"Redis Cache","type":"Redis","version":"7.0","memory":"8GB","persistence":"RDB","cluster_mode":true}'
redis-cli -p 6379 NODE.CREATE microservices session-cache cache '{"name":"Session Cache","type":"Redis","version":"7.0","memory":"4GB","ttl":"24h","persistence":"none"}'
redis-cli -p 6379 NODE.CREATE microservices product-cache cache '{"name":"Product Cache","type":"Redis","version":"7.0","memory":"16GB","eviction":"allkeys-lru","persistence":"AOF"}'

echo "Creating message queues..."
redis-cli -p 6379 NODE.CREATE microservices order-queue queue '{"name":"Order Queue","type":"RabbitMQ","version":"3.12","durability":true,"routing":"topic","max_length":10000}'
redis-cli -p 6379 NODE.CREATE microservices notification-queue queue '{"name":"Notification Queue","type":"Apache Kafka","version":"3.5","partitions":12,"replication_factor":3,"retention":"7d"}'
redis-cli -p 6379 NODE.CREATE microservices analytics-queue queue '{"name":"Analytics Queue","type":"Apache Kafka","version":"3.5","partitions":24,"replication_factor":3,"retention":"30d"}'

echo "Creating external services..."
redis-cli -p 6379 NODE.CREATE microservices payment-gateway external '{"name":"Stripe Payment Gateway","type":"external_api","endpoint":"https://api.stripe.com","rate_limit":"100/min","sla":"99.9%"}'
redis-cli -p 6379 NODE.CREATE microservices email-service external '{"name":"SendGrid Email Service","type":"external_api","endpoint":"https://api.sendgrid.com","rate_limit":"1000/min","sla":"99.95%"}'
redis-cli -p 6379 NODE.CREATE microservices sms-service external '{"name":"Twilio SMS Service","type":"external_api","endpoint":"https://api.twilio.com","rate_limit":"500/min","cost_per_message":"$0.0075"}'

echo "Creating monitoring and observability..."
redis-cli -p 6379 NODE.CREATE microservices prometheus monitoring '{"name":"Prometheus","type":"monitoring","version":"2.45","retention":"30d","scrape_interval":"15s","storage":"500GB"}'
redis-cli -p 6379 NODE.CREATE microservices grafana monitoring '{"name":"Grafana","type":"visualization","version":"10.0","dashboards":25,"alerts":50,"users":100}'
redis-cli -p 6379 NODE.CREATE microservices jaeger monitoring '{"name":"Jaeger","type":"tracing","version":"1.47","sampling_rate":"0.1","retention":"7d","storage":"elasticsearch"}'

echo "Creating service-to-database connections..."
redis-cli -p 6379 EDGE.CREATE microservices user-to-db user-service user-db connects '{"type":"database_connection","protocol":"TCP","pool_size":20,"timeout":"30s","ssl":true}'
redis-cli -p 6379 EDGE.CREATE microservices product-to-db product-service product-db connects '{"type":"database_connection","protocol":"TCP","pool_size":15,"timeout":"10s","read_preference":"secondary"}'
redis-cli -p 6379 EDGE.CREATE microservices order-to-db order-service order-db connects '{"type":"database_connection","protocol":"TCP","pool_size":25,"timeout":"30s","transaction_isolation":"read_committed"}'
redis-cli -p 6379 EDGE.CREATE microservices payment-to-db payment-service payment-db connects '{"type":"database_connection","protocol":"TCP","pool_size":10,"timeout":"60s","encryption":"TLS1.3"}'

echo "Creating cache connections..."
redis-cli -p 6379 EDGE.CREATE microservices user-to-session user-service session-cache uses '{"type":"session_storage","protocol":"Redis","ttl":"24h","serialization":"json"}'
redis-cli -p 6379 EDGE.CREATE microservices product-to-cache product-service product-cache uses '{"type":"cache_layer","protocol":"Redis","ttl":"1h","hit_ratio":"85%"}'
redis-cli -p 6379 EDGE.CREATE microservices order-to-cache order-service redis-cache uses '{"type":"cache_layer","protocol":"Redis","ttl":"30m","pattern":"write-through"}'

echo "Creating API Gateway routes..."
redis-cli -p 6379 EDGE.CREATE microservices gateway-to-user api-gateway user-service routes '{"path":"/api/users/*","method":"*","rate_limit":"1000/min","auth_required":true}'
redis-cli -p 6379 EDGE.CREATE microservices gateway-to-auth api-gateway auth-service routes '{"path":"/api/auth/*","method":"POST","rate_limit":"100/min","cors":true}'
redis-cli -p 6379 EDGE.CREATE microservices gateway-to-product api-gateway product-service routes '{"path":"/api/products/*","method":"*","rate_limit":"2000/min","cache":"5m"}'
redis-cli -p 6379 EDGE.CREATE microservices gateway-to-order api-gateway order-service routes '{"path":"/api/orders/*","method":"*","rate_limit":"500/min","auth_required":true}'
redis-cli -p 6379 EDGE.CREATE microservices gateway-to-payment api-gateway payment-service routes '{"path":"/api/payments/*","method":"POST","rate_limit":"200/min","auth_required":true,"encryption":true}'

echo "Creating service-to-service communications..."
redis-cli -p 6379 EDGE.CREATE microservices order-to-user order-service user-service calls '{"type":"sync_call","protocol":"HTTP","timeout":"5s","circuit_breaker":true,"retry_count":3}'
redis-cli -p 6379 EDGE.CREATE microservices order-to-product order-service product-service calls '{"type":"sync_call","protocol":"HTTP","timeout":"3s","cache":"30s","fallback":"cached_data"}'
redis-cli -p 6379 EDGE.CREATE microservices order-to-inventory order-service inventory-service calls '{"type":"sync_call","protocol":"HTTP","timeout":"2s","critical":true,"sla":"99.9%"}'
redis-cli -p 6379 EDGE.CREATE microservices order-to-payment order-service payment-service calls '{"type":"sync_call","protocol":"HTTPS","timeout":"30s","idempotent":true,"encryption":"end-to-end"}'

echo "Creating queue connections..."
redis-cli -p 6379 EDGE.CREATE microservices order-to-queue order-service order-queue publishes '{"event":"order_created","routing_key":"order.created","delivery_mode":"persistent","priority":5}'
redis-cli -p 6379 EDGE.CREATE microservices payment-to-queue payment-service order-queue publishes '{"event":"payment_processed","routing_key":"payment.processed","delivery_mode":"persistent","priority":8}'
redis-cli -p 6379 EDGE.CREATE microservices queue-to-notification order-queue notification-service consumes '{"event":"order_created","consumer_group":"notification_workers","batch_size":10,"max_wait":"5s"}'
redis-cli -p 6379 EDGE.CREATE microservices queue-to-analytics order-queue analytics-db streams '{"topic":"order_events","partition_key":"user_id","compression":"gzip","batch_size":1000}'

echo "Creating external service connections..."
redis-cli -p 6379 EDGE.CREATE microservices payment-to-stripe payment-service payment-gateway integrates '{"type":"payment_processing","webhook_url":"/webhooks/stripe","api_version":"2023-10-16","timeout":"30s"}'
redis-cli -p 6379 EDGE.CREATE microservices notification-to-email notification-service email-service integrates '{"type":"email_delivery","template_engine":"handlebars","tracking":true,"bounce_handling":true}'
redis-cli -p 6379 EDGE.CREATE microservices notification-to-sms notification-service sms-service integrates '{"type":"sms_delivery","country_code_validation":true,"delivery_reports":true,"opt_out_handling":true}'

echo "Creating monitoring connections..."
redis-cli -p 6379 EDGE.CREATE microservices user-to-prometheus user-service prometheus monitors '{"metrics":["http_requests_total","response_time","error_rate"],"scrape_interval":"15s","port":9090}'
redis-cli -p 6379 EDGE.CREATE microservices order-to-prometheus order-service prometheus monitors '{"metrics":["orders_created","orders_failed","processing_time"],"scrape_interval":"15s","port":9091}'
redis-cli -p 6379 EDGE.CREATE microservices payment-to-prometheus payment-service prometheus monitors '{"metrics":["payments_processed","payment_failures","fraud_detected"],"scrape_interval":"10s","port":9092}'

echo "✅ Microservices architecture created successfully!"
echo "📊 Graph contains:"
echo "   - 1 API Gateway with routing rules"
echo "   - 7 Microservices with detailed configurations"
echo "   - 5 Databases with different technologies"
echo "   - 3 Cache layers for performance"
echo "   - 3 Message queues for async processing"
echo "   - 3 External service integrations"
echo "   - 3 Monitoring and observability tools"
echo "   - 25+ edges showing complex relationships"
