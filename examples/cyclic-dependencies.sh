#!/bin/bash

# Cyclic Dependencies Sample Data
# Creates intentional cyclic patterns to showcase cycle detection capabilities

echo "🔄 Creating cyclic dependencies graph..."

# Create the main graph
redis-cli -p 6379 GRAPH.CREATE cyclic_deps "Cyclic dependency patterns and circular references"

echo "Creating software modules with circular dependencies..."
redis-cli -p 6379 NODE.CREATE cyclic_deps module-auth module '{"name":"Authentication Module","language":"Java","version":"2.1.0","dependencies":["user-mgmt","session-mgmt"],"circular_risk":"high","refactor_priority":"urgent"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps module-user-mgmt module '{"name":"User Management Module","language":"Java","version":"1.8.0","dependencies":["auth","profile-mgmt"],"circular_risk":"high","refactor_priority":"urgent"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps module-session-mgmt module '{"name":"Session Management Module","language":"Java","version":"1.5.0","dependencies":["auth","user-mgmt"],"circular_risk":"medium","refactor_priority":"high"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps module-profile-mgmt module '{"name":"Profile Management Module","language":"Java","version":"2.0.0","dependencies":["user-mgmt","notification"],"circular_risk":"low","refactor_priority":"medium"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps module-notification module '{"name":"Notification Module","language":"Java","version":"1.3.0","dependencies":["user-mgmt","email-service"],"circular_risk":"low","refactor_priority":"low"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps module-email-service module '{"name":"Email Service Module","language":"Java","version":"1.1.0","dependencies":["template-engine","user-mgmt"],"circular_risk":"medium","refactor_priority":"medium"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps module-template-engine module '{"name":"Template Engine Module","language":"Java","version":"1.0.0","dependencies":["profile-mgmt"],"circular_risk":"low","refactor_priority":"low"}'

echo "Creating database tables with foreign key cycles..."
redis-cli -p 6379 NODE.CREATE cyclic_deps table-users table '{"name":"users","primary_key":"user_id","foreign_keys":["default_address_id","manager_id"],"circular_refs":["addresses","employees"],"constraint_issues":"deferred"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps table-addresses table '{"name":"addresses","primary_key":"address_id","foreign_keys":["user_id","country_id"],"circular_refs":["users"],"constraint_issues":"none"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps table-employees table '{"name":"employees","primary_key":"employee_id","foreign_keys":["user_id","department_id","supervisor_id"],"circular_refs":["departments","self_reference"],"constraint_issues":"self_referential"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps table-departments table '{"name":"departments","primary_key":"dept_id","foreign_keys":["head_employee_id","parent_dept_id"],"circular_refs":["employees","self_reference"],"constraint_issues":"hierarchical_cycle"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps table-projects table '{"name":"projects","primary_key":"project_id","foreign_keys":["lead_employee_id","parent_project_id"],"circular_refs":["employees","project_dependencies"],"constraint_issues":"dependency_cycle"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps table-project-deps table '{"name":"project_dependencies","primary_key":"dependency_id","foreign_keys":["project_id","depends_on_project_id"],"circular_refs":["projects"],"constraint_issues":"circular_dependencies"}'

echo "Creating microservices with circular communication..."
redis-cli -p 6379 NODE.CREATE cyclic_deps service-order microservice '{"name":"Order Service","port":8001,"dependencies":["payment-service","inventory-service"],"calls":["payment","inventory","customer"],"circular_calls":"payment->customer->order"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps service-payment microservice '{"name":"Payment Service","port":8002,"dependencies":["order-service","customer-service"],"calls":["order","customer","fraud"],"circular_calls":"order->payment->customer->order"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps service-customer microservice '{"name":"Customer Service","port":8003,"dependencies":["order-service","payment-service"],"calls":["order","payment","notification"],"circular_calls":"payment->customer->order->payment"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps service-inventory microservice '{"name":"Inventory Service","port":8004,"dependencies":["order-service","supplier-service"],"calls":["order","supplier"],"circular_calls":"none_direct"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps service-supplier microservice '{"name":"Supplier Service","port":8005,"dependencies":["inventory-service","order-service"],"calls":["inventory","order"],"circular_calls":"inventory->supplier->order->inventory"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps service-fraud microservice '{"name":"Fraud Detection Service","port":8006,"dependencies":["payment-service","customer-service"],"calls":["payment","customer"],"circular_calls":"none"}'

echo "Creating workflow processes with circular flows..."
redis-cli -p 6379 NODE.CREATE cyclic_deps process-approval workflow '{"name":"Approval Process","type":"business_process","steps":["submit","review","approve","notify"],"circular_flow":"reject->resubmit->review","max_iterations":"3"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps process-review workflow '{"name":"Review Process","type":"business_process","steps":["assign","evaluate","feedback","decision"],"circular_flow":"feedback->revise->evaluate","max_iterations":"5"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps process-revision workflow '{"name":"Revision Process","type":"business_process","steps":["receive_feedback","modify","resubmit","validate"],"circular_flow":"validate->feedback->modify","max_iterations":"unlimited"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps process-validation workflow '{"name":"Validation Process","type":"business_process","steps":["check","test","report","approve"],"circular_flow":"fail->revision->validation","max_iterations":"10"}'

echo "Creating organizational hierarchy with circular reporting..."
redis-cli -p 6379 NODE.CREATE cyclic_deps manager-alice person '{"name":"Alice Johnson","role":"Engineering Manager","reports_to":"bob_smith","manages":["charlie_brown","diana_prince"],"circular_issue":"reports_to_subordinate"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps manager-bob person '{"name":"Bob Smith","role":"Technical Lead","reports_to":"charlie_brown","manages":["alice_johnson"],"circular_issue":"manages_own_manager"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps manager-charlie person '{"name":"Charlie Brown","role":"Senior Developer","reports_to":"alice_johnson","manages":["bob_smith"],"circular_issue":"circular_reporting"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps manager-diana person '{"name":"Diana Prince","role":"Product Manager","reports_to":"alice_johnson","manages":["eve_wilson"],"circular_issue":"none"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps manager-eve person '{"name":"Eve Wilson","role":"UX Designer","reports_to":"diana_prince","manages":["alice_johnson"],"circular_issue":"skip_level_cycle"}'

echo "Creating network topology with routing loops..."
redis-cli -p 6379 NODE.CREATE cyclic_deps router-a network '{"name":"Router A","ip":"192.168.1.1","routes_to":["router_b","router_c"],"routing_protocol":"OSPF","loop_prevention":"split_horizon","metric":"10"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps router-b network '{"name":"Router B","ip":"192.168.1.2","routes_to":["router_c","router_d"],"routing_protocol":"OSPF","loop_prevention":"poison_reverse","metric":"15"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps router-c network '{"name":"Router C","ip":"192.168.1.3","routes_to":["router_d","router_a"],"routing_protocol":"OSPF","loop_prevention":"hold_down","metric":"20"}'
redis-cli -p 6379 NODE.CREATE cyclic_deps router-d network '{"name":"Router D","ip":"192.168.1.4","routes_to":["router_a","router_b"],"routing_protocol":"OSPF","loop_prevention":"ttl_decrement","metric":"12"}'

echo "Creating module dependency cycles (3-node cycle)..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps auth-to-user module-auth module-user-mgmt depends_on '{"type":"compile_time","severity":"critical","cycle_length":"3","refactor_effort":"high","breaking_change":"true"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps user-to-session module-user-mgmt module-session-mgmt depends_on '{"type":"runtime","severity":"high","cycle_length":"3","refactor_effort":"medium","breaking_change":"false"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps session-to-auth module-session-mgmt module-auth depends_on '{"type":"compile_time","severity":"critical","cycle_length":"3","refactor_effort":"high","breaking_change":"true"}'

echo "Creating additional module cycles (4-node cycle)..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps user-to-profile module-user-mgmt module-profile-mgmt depends_on '{"type":"runtime","severity":"medium","cycle_length":"4","refactor_effort":"low","breaking_change":"false"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps profile-to-template module-profile-mgmt module-template-engine depends_on '{"type":"compile_time","severity":"low","cycle_length":"4","refactor_effort":"low","breaking_change":"false"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps template-to-email module-template-engine module-email-service depends_on '{"type":"runtime","severity":"medium","cycle_length":"4","refactor_effort":"medium","breaking_change":"false"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps email-to-user module-email-service module-user-mgmt depends_on '{"type":"runtime","severity":"medium","cycle_length":"4","refactor_effort":"medium","breaking_change":"false"}'

echo "Creating database foreign key cycles..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps users-to-addresses table-users table-addresses references '{"foreign_key":"default_address_id","constraint":"deferred","nullable":"true","cycle_type":"optional_reference","resolution":"nullable_fk"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps addresses-to-users table-addresses table-users references '{"foreign_key":"user_id","constraint":"immediate","nullable":"false","cycle_type":"required_reference","resolution":"creation_order"}'

echo "Creating employee hierarchy cycles..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps employees-to-dept table-employees table-departments belongs_to '{"foreign_key":"department_id","constraint":"immediate","nullable":"false","cycle_type":"hierarchical","resolution":"department_first"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps dept-to-head table-departments table-employees has_head '{"foreign_key":"head_employee_id","constraint":"deferred","nullable":"true","cycle_type":"leadership","resolution":"nullable_head"}'

echo "Creating self-referential cycles..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps employee-to-supervisor table-employees table-employees reports_to '{"foreign_key":"supervisor_id","constraint":"immediate","nullable":"true","cycle_type":"self_referential","resolution":"hierarchy_validation"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps dept-to-parent table-departments table-departments child_of '{"foreign_key":"parent_dept_id","constraint":"immediate","nullable":"true","cycle_type":"tree_structure","resolution":"depth_limit"}'

echo "Creating project dependency cycles..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps project-to-deps table-projects table-project-deps has_dependencies '{"foreign_key":"project_id","constraint":"immediate","nullable":"false","cycle_type":"dependency_graph","resolution":"topological_sort"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps deps-to-project table-project-deps table-projects depends_on '{"foreign_key":"depends_on_project_id","constraint":"immediate","nullable":"false","cycle_type":"circular_dependency","resolution":"cycle_detection"}'

echo "Creating microservice communication cycles..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps order-to-payment service-order service-payment calls '{"protocol":"HTTP","timeout":"30s","circuit_breaker":"enabled","cycle_risk":"high","async_option":"available"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps payment-to-customer service-payment service-customer calls '{"protocol":"HTTP","timeout":"15s","circuit_breaker":"enabled","cycle_risk":"high","async_option":"recommended"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps customer-to-order service-customer service-order calls '{"protocol":"HTTP","timeout":"20s","circuit_breaker":"enabled","cycle_risk":"critical","async_option":"required"}'

echo "Creating inventory-supplier cycle..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps inventory-to-supplier service-inventory service-supplier requests '{"protocol":"HTTP","timeout":"60s","circuit_breaker":"enabled","cycle_risk":"medium","async_option":"preferred"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps supplier-to-inventory service-supplier service-inventory updates '{"protocol":"HTTP","timeout":"45s","circuit_breaker":"enabled","cycle_risk":"medium","async_option":"preferred"}'

echo "Creating workflow process cycles..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps approval-to-review process-approval process-review triggers '{"condition":"requires_review","max_iterations":"3","cycle_detection":"enabled","timeout":"24h","escalation":"manager"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps review-to-revision process-review process-revision requests '{"condition":"needs_changes","max_iterations":"5","cycle_detection":"enabled","timeout":"48h","escalation":"senior_reviewer"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps revision-to-validation process-revision process-validation submits_to '{"condition":"changes_made","max_iterations":"unlimited","cycle_detection":"warning","timeout":"72h","escalation":"none"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps validation-to-review process-validation process-review returns_to '{"condition":"validation_failed","max_iterations":"10","cycle_detection":"error","timeout":"12h","escalation":"team_lead"}'

echo "Creating organizational reporting cycles..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps alice-to-bob manager-alice manager-bob reports_to '{"relationship":"manager_to_lead","cycle_issue":"manages_own_manager","severity":"critical","resolution_required":"immediate"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps bob-to-charlie manager-bob manager-charlie reports_to '{"relationship":"lead_to_developer","cycle_issue":"circular_reporting","severity":"critical","resolution_required":"immediate"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps charlie-to-alice manager-charlie manager-alice reports_to '{"relationship":"developer_to_manager","cycle_issue":"completes_cycle","severity":"critical","resolution_required":"immediate"}'

echo "Creating skip-level reporting cycle..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps alice-to-diana manager-alice manager-diana manages '{"relationship":"manager_to_pm","cycle_issue":"none","severity":"none","resolution_required":"none"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps diana-to-eve manager-diana manager-eve manages '{"relationship":"pm_to_designer","cycle_issue":"none","severity":"none","resolution_required":"none"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps eve-to-alice manager-eve manager-alice influences '{"relationship":"designer_to_manager","cycle_issue":"skip_level_cycle","severity":"medium","resolution_required":"policy_clarification"}'

echo "Creating network routing loops..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps router-a-to-b router-a router-b routes_to '{"protocol":"OSPF","metric":"10","next_hop":"direct","loop_prevention":"split_horizon","backup_route":"router_c"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps router-b-to-c router-b router-c routes_to '{"protocol":"OSPF","metric":"15","next_hop":"direct","loop_prevention":"poison_reverse","backup_route":"router_d"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps router-c-to-d router-c router-d routes_to '{"protocol":"OSPF","metric":"20","next_hop":"direct","loop_prevention":"hold_down","backup_route":"router_a"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps router-d-to-a router-d router-a routes_to '{"protocol":"OSPF","metric":"12","next_hop":"direct","loop_prevention":"ttl_decrement","backup_route":"router_b"}'

echo "Creating additional routing paths for complex cycles..."
redis-cli -p 6379 EDGE.CREATE cyclic_deps router-a-to-c router-a router-c backup_route '{"protocol":"OSPF","metric":"25","next_hop":"multi_hop","loop_prevention":"path_vector","primary":"false"}'
redis-cli -p 6379 EDGE.CREATE cyclic_deps router-b-to-d router-b router-d backup_route '{"protocol":"OSPF","metric":"30","next_hop":"multi_hop","loop_prevention":"distance_vector","primary":"false"}'

echo "✅ Cyclic dependencies graph created successfully!"
echo "📊 Graph contains multiple cycle types:"
echo ""
echo "🔄 **3-Node Cycles:**"
echo "   - Module Dependencies: Auth → User Management → Session → Auth"
echo "   - Organizational: Alice → Bob → Charlie → Alice (reporting cycle)"
echo ""
echo "🔄 **4-Node Cycles:**"
echo "   - Extended Modules: User → Profile → Template → Email → User"
echo "   - Network Routing: Router A → B → C → D → A"
echo ""
echo "🔄 **2-Node Cycles:**"
echo "   - Database: Users ↔ Addresses (foreign key cycle)"
echo "   - Microservices: Inventory ↔ Supplier (communication cycle)"
echo ""
echo "🔄 **Self-Referential Cycles:**"
echo "   - Employee → Supervisor (same table)"
echo "   - Department → Parent Department (hierarchical)"
echo ""
echo "🔄 **Complex Multi-Path Cycles:**"
echo "   - Microservices: Order → Payment → Customer → Order"
echo "   - Workflows: Approval → Review → Revision → Validation → Review"
echo ""
echo "🔄 **Cycle Detection Features Demonstrated:**"
echo "   - Different cycle lengths (2, 3, 4+ nodes)"
echo "   - Various relationship types (dependencies, references, calls)"
echo "   - Self-referential patterns"
echo "   - Skip-level and hierarchical cycles"
echo "   - Network routing loops with backup paths"
echo "   - Business process circular flows"
echo ""
echo "⚠️  **Cycle Severity Levels:**"
echo "   - Critical: Compile-time dependencies, organizational reporting"
echo "   - High: Runtime dependencies, service communication"
echo "   - Medium: Optional references, workflow iterations"
echo "   - Low: Template dependencies, backup routes"
