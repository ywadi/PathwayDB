#!/bin/bash

# Zero-Code Workflow Automation Sample Data
# Creates complex workflow automation scenarios with triggers, conditions, and actions

echo "🔄 Creating zero-code workflow automation graph..."

# Create the main graph
redis-cli -p 6379 GRAPH.CREATE workflows "Zero-code workflow automation platform"

echo "Creating workflow triggers..."
redis-cli -p 6379 NODE.CREATE workflows webhook-trigger trigger '{"name":"Webhook Trigger","type":"http_webhook","endpoint":"/webhook/customer-signup","method":"POST","authentication":"api_key","rate_limit":"1000/hour"}'
redis-cli -p 6379 NODE.CREATE workflows schedule-trigger trigger '{"name":"Daily Report Trigger","type":"cron_schedule","expression":"0 9 * * MON-FRI","timezone":"UTC","enabled":true,"next_run":"2024-01-15T09:00:00Z"}'
redis-cli -p 6379 NODE.CREATE workflows email-trigger trigger '{"name":"Email Received Trigger","type":"email_webhook","provider":"gmail","folder":"inbox","filter":"from:support@company.com","polling_interval":"5m"}'
redis-cli -p 6379 NODE.CREATE workflows file-trigger trigger '{"name":"File Upload Trigger","type":"file_watcher","path":"/uploads/invoices/","extensions":[".pdf",".xlsx"],"action":"created","recursive":true}'
redis-cli -p 6379 NODE.CREATE workflows database-trigger trigger '{"name":"New Order Trigger","type":"database_change","table":"orders","operation":"INSERT","conditions":"status=pending","debounce":"30s"}'

echo "Creating condition nodes..."
redis-cli -p 6379 NODE.CREATE workflows amount-condition condition '{"name":"Order Amount Check","type":"numeric_comparison","field":"order.total","operator":"greater_than","value":1000,"currency":"USD"}'
redis-cli -p 6379 NODE.CREATE workflows region-condition condition '{"name":"Customer Region Check","type":"text_match","field":"customer.region","operator":"equals","value":"North America","case_sensitive":false}'
redis-cli -p 6379 NODE.CREATE workflows time-condition condition '{"name":"Business Hours Check","type":"time_range","start_time":"09:00","end_time":"17:00","timezone":"EST","weekdays_only":true}'
redis-cli -p 6379 NODE.CREATE workflows inventory-condition condition '{"name":"Stock Level Check","type":"numeric_comparison","field":"product.stock","operator":"less_than","value":10,"alert_threshold":5}'
redis-cli -p 6379 NODE.CREATE workflows approval-condition condition '{"name":"Manager Approval Required","type":"boolean_check","field":"order.requires_approval","value":true,"timeout":"24h","escalation":true}'

echo "Creating data transformation nodes..."
redis-cli -p 6379 NODE.CREATE workflows data-mapper transformer '{"name":"Customer Data Mapper","type":"field_mapping","mappings":{"email":"customer_email","name":"full_name","phone":"contact_number"},"validation":true,"required_fields":["email","name"]}'
redis-cli -p 6379 NODE.CREATE workflows json-parser transformer '{"name":"JSON Parser","type":"data_parser","input_format":"json","output_format":"structured","schema_validation":true,"error_handling":"skip_invalid"}'
redis-cli -p 6379 NODE.CREATE workflows date-formatter transformer '{"name":"Date Formatter","type":"date_transform","input_format":"ISO8601","output_format":"MM/DD/YYYY","timezone_conversion":"UTC_to_local"}'
redis-cli -p 6379 NODE.CREATE workflows text-processor transformer '{"name":"Text Processor","type":"text_transform","operations":["trim","lowercase","remove_special_chars"],"encoding":"UTF-8"}'
redis-cli -p 6379 NODE.CREATE workflows aggregator transformer '{"name":"Order Aggregator","type":"data_aggregation","group_by":"customer_id","functions":["sum(amount)","count(*)","avg(rating)"],"window":"24h"}'

echo "Creating action nodes..."
redis-cli -p 6379 NODE.CREATE workflows email-action action '{"name":"Send Welcome Email","type":"email_send","template":"welcome_template","provider":"sendgrid","personalization":true,"tracking":true,"retry_count":3}'
redis-cli -p 6379 NODE.CREATE workflows slack-action action '{"name":"Notify Sales Team","type":"slack_message","channel":"#sales-alerts","template":"New high-value order: {{order.id}} - ${{order.total}}","mention_users":["@sales-manager"]}'
redis-cli -p 6379 NODE.CREATE workflows database-action action '{"name":"Update Customer Status","type":"database_update","table":"customers","set_fields":{"status":"premium","updated_at":"{{now}}"},"where_condition":"id={{customer.id}}"}'
redis-cli -p 6379 NODE.CREATE workflows api-action action '{"name":"Create CRM Contact","type":"http_request","method":"POST","url":"https://api.crm.com/contacts","headers":{"Authorization":"Bearer {{crm_token}}"},"body_template":"crm_contact.json"}'
redis-cli -p 6379 NODE.CREATE workflows file-action action '{"name":"Generate Invoice PDF","type":"pdf_generation","template":"invoice_template.html","output_path":"/invoices/{{order.id}}.pdf","watermark":true}'
redis-cli -p 6379 NODE.CREATE workflows sms-action action '{"name":"Send Order Confirmation","type":"sms_send","provider":"twilio","template":"Order {{order.id}} confirmed. Total: ${{order.total}}","country_code_validation":true}'

echo "Creating integration nodes..."
redis-cli -p 6379 NODE.CREATE workflows salesforce-integration integration '{"name":"Salesforce CRM","type":"salesforce","instance":"https://company.salesforce.com","api_version":"v58.0","oauth_flow":"server_to_server","sync_frequency":"15m"}'
redis-cli -p 6379 NODE.CREATE workflows hubspot-integration integration '{"name":"HubSpot Marketing","type":"hubspot","portal_id":"12345678","api_key_encrypted":true,"contact_sync":true,"deal_sync":true,"webhook_validation":true}'
redis-cli -p 6379 NODE.CREATE workflows stripe-integration integration '{"name":"Stripe Payments","type":"stripe","webhook_endpoint":"/webhooks/stripe","events":["payment_intent.succeeded","invoice.payment_failed"],"signature_verification":true}'
redis-cli -p 6379 NODE.CREATE workflows zapier-integration integration '{"name":"Zapier Connector","type":"zapier","webhook_url":"https://hooks.zapier.com/hooks/catch/12345/abcdef","authentication":"api_key","rate_limit":"100/min"}'

echo "Creating workflow orchestration nodes..."
redis-cli -p 6379 NODE.CREATE workflows parallel-processor orchestrator '{"name":"Parallel Task Processor","type":"parallel_execution","max_concurrent":5,"timeout":"300s","failure_strategy":"continue_on_error","result_aggregation":true}'
redis-cli -p 6379 NODE.CREATE workflows sequential-processor orchestrator '{"name":"Sequential Task Processor","type":"sequential_execution","stop_on_error":true,"checkpoint_frequency":"every_step","rollback_support":true}'
redis-cli -p 6379 NODE.CREATE workflows loop-processor orchestrator '{"name":"Batch Loop Processor","type":"loop_execution","batch_size":100,"delay_between_batches":"5s","progress_tracking":true,"resume_on_failure":true}'
redis-cli -p 6379 NODE.CREATE workflows decision-router orchestrator '{"name":"Decision Router","type":"conditional_routing","default_path":"fallback","evaluation_order":"priority","caching":true,"audit_trail":true}'

echo "Creating error handling nodes..."
redis-cli -p 6379 NODE.CREATE workflows retry-handler error_handler '{"name":"Retry Handler","type":"retry_logic","max_attempts":3,"backoff_strategy":"exponential","base_delay":"1s","max_delay":"60s","retry_conditions":["timeout","rate_limit"]}'
redis-cli -p 6379 NODE.CREATE workflows dead-letter-queue error_handler '{"name":"Dead Letter Queue","type":"failed_message_queue","retention":"7d","max_size":"10000","alert_threshold":"100","reprocessing":true}'
redis-cli -p 6379 NODE.CREATE workflows error-notification error_handler '{"name":"Error Alert System","type":"error_notification","channels":["email","slack","pagerduty"],"severity_levels":["low","medium","high","critical"],"escalation_rules":true}'

echo "Creating workflow connections - Customer Onboarding Flow..."
redis-cli -p 6379 EDGE.CREATE workflows webhook-to-mapper webhook-trigger data-mapper triggers '{"event":"customer_signup","payload_validation":true,"rate_limit_check":true}'
redis-cli -p 6379 EDGE.CREATE workflows mapper-to-region data-mapper region-condition processes '{"field_extraction":"customer.region","validation_rules":["required","valid_region"]}'
redis-cli -p 6379 EDGE.CREATE workflows region-to-email region-condition email-action executes '{"condition_result":"true","template_selection":"region_specific","personalization_data":"customer_profile"}'
redis-cli -p 6379 EDGE.CREATE workflows region-to-crm region-condition api-action executes '{"condition_result":"true","data_mapping":"crm_contact_format","duplicate_check":true}'

echo "Creating workflow connections - Order Processing Flow..."
redis-cli -p 6379 EDGE.CREATE workflows database-to-amount database-trigger amount-condition evaluates '{"trigger_data":"order_details","condition_field":"total_amount","currency_conversion":true}'
redis-cli -p 6379 EDGE.CREATE workflows amount-to-approval amount-condition approval-condition routes '{"high_value_threshold":1000,"approval_workflow":"manager_approval","timeout":"24h"}'
redis-cli -p 6379 EDGE.CREATE workflows approval-to-slack approval-condition slack-action notifies '{"approval_required":"true","urgency":"high","escalation_timer":"2h"}'
redis-cli -p 6379 EDGE.CREATE workflows amount-to-invoice amount-condition file-action generates '{"condition_result":"true","template":"premium_invoice","digital_signature":true}'

echo "Creating workflow connections - Scheduled Reporting Flow..."
redis-cli -p 6379 EDGE.CREATE workflows schedule-to-aggregator schedule-trigger aggregator triggers '{"schedule_type":"daily","data_range":"24h","timezone_handling":true}'
redis-cli -p 6379 EDGE.CREATE workflows aggregator-to-time aggregator time-condition checks '{"business_hours_only":true,"holiday_calendar":"US","skip_weekends":true}'
redis-cli -p 6379 EDGE.CREATE workflows time-to-email time-condition email-action sends '{"condition_result":"true","report_format":"html","attachment_support":true}'
redis-cli -p 6379 EDGE.CREATE workflows aggregator-to-database aggregator database-action stores '{"table":"daily_reports","data_retention":"1_year","compression":true}'

echo "Creating integration connections..."
redis-cli -p 6379 EDGE.CREATE workflows api-to-salesforce api-action salesforce-integration syncs '{"object_type":"Contact","field_mapping":"salesforce_schema","duplicate_handling":"merge"}'
redis-cli -p 6379 EDGE.CREATE workflows email-to-hubspot email-action hubspot-integration tracks '{"email_events":["sent","opened","clicked"],"contact_scoring":true,"campaign_attribution":true}'
redis-cli -p 6379 EDGE.CREATE workflows database-to-stripe database-action stripe-integration processes '{"payment_method":"automatic","invoice_generation":true,"webhook_confirmation":true}'

echo "Creating error handling connections..."
redis-cli -p 6379 EDGE.CREATE workflows email-to-retry email-action retry-handler handles '{"error_types":["timeout","rate_limit","temporary_failure"],"max_retries":3,"backoff_exponential":true}'
redis-cli -p 6379 EDGE.CREATE workflows retry-to-dlq retry-handler dead-letter-queue escalates '{"max_attempts_reached":true,"permanent_failure":true,"investigation_required":true}'
redis-cli -p 6379 EDGE.CREATE workflows dlq-to-alert dead-letter-queue error-notification triggers '{"threshold":"10_messages","time_window":"1h","severity":"high"}'

echo "Creating orchestration connections..."
redis-cli -p 6379 EDGE.CREATE workflows mapper-to-parallel data-mapper parallel-processor distributes '{"task_distribution":"round_robin","load_balancing":true,"resource_monitoring":true}'
redis-cli -p 6379 EDGE.CREATE workflows parallel-to-decision parallel-processor decision-router aggregates '{"result_collection":"all_tasks","success_criteria":"80%","failure_threshold":"20%"}'
redis-cli -p 6379 EDGE.CREATE workflows decision-to-sequential decision-router sequential-processor routes '{"condition":"success_rate_high","checkpoint_enabled":true,"rollback_plan":true}'

echo "✅ Zero-code workflow automation graph created successfully!"
echo "📊 Graph contains:"
echo "   - 5 Different trigger types (webhook, schedule, email, file, database)"
echo "   - 5 Conditional logic nodes with complex rules"
echo "   - 5 Data transformation and processing nodes"
echo "   - 6 Action nodes for various integrations"
echo "   - 4 Third-party service integrations"
echo "   - 4 Workflow orchestration patterns"
echo "   - 3 Error handling and recovery mechanisms"
echo "   - 20+ edges showing complex workflow relationships"
echo ""
echo "🔄 Workflow Examples Created:"
echo "   1. Customer Onboarding: Webhook → Data Mapping → Region Check → Email + CRM"
echo "   2. Order Processing: Database Change → Amount Check → Approval → Notifications"
echo "   3. Scheduled Reporting: Cron → Data Aggregation → Business Hours → Email Report"
echo "   4. Error Handling: Failed Actions → Retry Logic → Dead Letter Queue → Alerts"
