#!/bin/bash

# Sample data script for PathwayDB IDE testing
echo "🔧 Creating sample data for PathwayDB IDE..."

# Connect to Redis and create sample graph data
echo "Creating microservices graph..."
redis-cli -p 6379 GRAPH.CREATE microservices "Sample microservices architecture"

echo "Creating service nodes..."
redis-cli -p 6379 NODE.CREATE microservices user-service service '{"name":"User Service","port":8001,"language":"Go","replicas":3}'
redis-cli -p 6379 NODE.CREATE microservices order-service service '{"name":"Order Service","port":8002,"language":"Python","replicas":2}'
redis-cli -p 6379 NODE.CREATE microservices payment-service service '{"name":"Payment Service","port":8003,"language":"Java","replicas":2}'
redis-cli -p 6379 NODE.CREATE microservices notification-service service '{"name":"Notification Service","port":8004,"language":"Node.js","replicas":1}'

echo "Creating database nodes..."
redis-cli -p 6379 NODE.CREATE microservices user-db database '{"name":"User Database","type":"PostgreSQL","version":"14.0"}'
redis-cli -p 6379 NODE.CREATE microservices order-db database '{"name":"Order Database","type":"MongoDB","version":"5.0"}'
redis-cli -p 6379 NODE.CREATE microservices payment-db database '{"name":"Payment Database","type":"MySQL","version":"8.0"}'

echo "Creating cache nodes..."
redis-cli -p 6379 NODE.CREATE microservices redis-cache cache '{"name":"Redis Cache","type":"Redis","version":"7.0","memory":"2GB"}'
redis-cli -p 6379 NODE.CREATE microservices session-cache cache '{"name":"Session Cache","type":"Redis","version":"7.0","memory":"1GB"}'

echo "Creating API gateway..."
redis-cli -p 6379 NODE.CREATE microservices api-gateway api '{"name":"API Gateway","type":"Kong","version":"3.0"}'

echo "Creating service dependencies..."
redis-cli -p 6379 EDGE.CREATE microservices svc-user-db user-service user-db connects '{"type":"database_connection","protocol":"TCP"}'
redis-cli -p 6379 EDGE.CREATE microservices svc-order-db order-service order-db connects '{"type":"database_connection","protocol":"TCP"}'
redis-cli -p 6379 EDGE.CREATE microservices svc-payment-db payment-service payment-db connects '{"type":"database_connection","protocol":"TCP"}'

redis-cli -p 6379 EDGE.CREATE microservices svc-user-cache user-service redis-cache uses '{"type":"cache_connection","protocol":"Redis"}'
redis-cli -p 6379 EDGE.CREATE microservices svc-session-cache user-service session-cache uses '{"type":"session_storage","protocol":"Redis"}'

redis-cli -p 6379 EDGE.CREATE microservices order-user order-service user-service calls '{"type":"http_api","method":"GET,POST"}'
redis-cli -p 6379 EDGE.CREATE microservices payment-order payment-service order-service calls '{"type":"http_api","method":"POST"}'
redis-cli -p 6379 EDGE.CREATE microservices notification-order notification-service order-service subscribes '{"type":"event_stream","protocol":"Kafka"}'
redis-cli -p 6379 EDGE.CREATE microservices notification-payment notification-service payment-service subscribes '{"type":"event_stream","protocol":"Kafka"}'

redis-cli -p 6379 EDGE.CREATE microservices gateway-user api-gateway user-service routes '{"type":"http_proxy","path":"/api/users"}'
redis-cli -p 6379 EDGE.CREATE microservices gateway-order api-gateway order-service routes '{"type":"http_proxy","path":"/api/orders"}'
redis-cli -p 6379 EDGE.CREATE microservices gateway-payment api-gateway payment-service routes '{"type":"http_proxy","path":"/api/payments"}'

echo "Creating e-commerce graph..."
redis-cli -p 6379 GRAPH.CREATE ecommerce "E-commerce platform architecture"

echo "Creating e-commerce nodes..."
redis-cli -p 6379 NODE.CREATE ecommerce frontend user '{"name":"Web Frontend","framework":"React","version":"18.0"}'
redis-cli -p 6379 NODE.CREATE ecommerce mobile-app user '{"name":"Mobile App","platform":"React Native","version":"0.72"}'
redis-cli -p 6379 NODE.CREATE ecommerce product-catalog service '{"name":"Product Catalog","language":"Go","port":9001}'
redis-cli -p 6379 NODE.CREATE ecommerce shopping-cart service '{"name":"Shopping Cart","language":"Python","port":9002}'
redis-cli -p 6379 NODE.CREATE ecommerce checkout service '{"name":"Checkout Service","language":"Java","port":9003}'
redis-cli -p 6379 NODE.CREATE ecommerce inventory service '{"name":"Inventory Service","language":"Go","port":9004}'

redis-cli -p 6379 NODE.CREATE ecommerce catalog-db database '{"name":"Catalog DB","type":"PostgreSQL","size":"100GB"}'
redis-cli -p 6379 NODE.CREATE ecommerce cart-db database '{"name":"Cart DB","type":"Redis","size":"10GB"}'
redis-cli -p 6379 NODE.CREATE ecommerce order-db database '{"name":"Order DB","type":"MongoDB","size":"50GB"}'

echo "Creating e-commerce edges..."
redis-cli -p 6379 EDGE.CREATE ecommerce frontend-catalog frontend product-catalog requests '{"type":"GraphQL","endpoint":"/graphql"}'
redis-cli -p 6379 EDGE.CREATE ecommerce frontend-cart frontend shopping-cart requests '{"type":"REST","endpoint":"/api/cart"}'
redis-cli -p 6379 EDGE.CREATE ecommerce mobile-catalog mobile-app product-catalog requests '{"type":"GraphQL","endpoint":"/graphql"}'
redis-cli -p 6379 EDGE.CREATE ecommerce mobile-cart mobile-app shopping-cart requests '{"type":"REST","endpoint":"/api/cart"}'

redis-cli -p 6379 EDGE.CREATE ecommerce catalog-db-conn product-catalog catalog-db connects '{"type":"SQL","pool_size":20}'
redis-cli -p 6379 EDGE.CREATE ecommerce cart-db-conn shopping-cart cart-db connects '{"type":"Redis","pool_size":10}'
redis-cli -p 6379 EDGE.CREATE ecommerce checkout-inventory checkout inventory checks '{"type":"gRPC","timeout":"5s"}'
redis-cli -p 6379 EDGE.CREATE ecommerce checkout-order checkout order-db stores '{"type":"MongoDB","collection":"orders"}'

echo "✅ Sample data created successfully!"
echo ""
echo "📊 Created graphs:"
echo "  • microservices - Sample microservices architecture (11 nodes, 12 edges)"
echo "  • ecommerce - E-commerce platform architecture (9 nodes, 8 edges)"
echo ""
echo "🚀 Start the IDE with: ./ide/start.sh"
