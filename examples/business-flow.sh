#!/bin/bash

# Business Process Flow Sample Data
# Creates a comprehensive business workflow with departments, processes, and dependencies

echo "🏢 Creating business process flow graph..."

# Create the main graph
redis-cli -p 6379 GRAPH.CREATE business_flow "Enterprise business process management"

echo "Creating organizational departments..."
redis-cli -p 6379 NODE.CREATE business_flow sales-dept department '{"name":"Sales Department","head":"Sarah Johnson","employees":25,"budget":"$2.5M","location":"Floor 3","kpis":["revenue","conversion_rate","pipeline_value"]}'
redis-cli -p 6379 NODE.CREATE business_flow marketing-dept department '{"name":"Marketing Department","head":"Mike Chen","employees":15,"budget":"$1.8M","location":"Floor 2","kpis":["lead_generation","brand_awareness","roi"]}'
redis-cli -p 6379 NODE.CREATE business_flow finance-dept department '{"name":"Finance Department","head":"Lisa Rodriguez","employees":12,"budget":"$1.2M","location":"Floor 4","kpis":["cash_flow","profit_margin","cost_control"]}'
redis-cli -p 6379 NODE.CREATE business_flow hr-dept department '{"name":"Human Resources","head":"David Kim","employees":8,"budget":"$800K","location":"Floor 1","kpis":["employee_satisfaction","retention_rate","time_to_hire"]}'
redis-cli -p 6379 NODE.CREATE business_flow operations-dept department '{"name":"Operations","head":"Jennifer Smith","employees":35,"budget":"$3.2M","location":"Floor 5","kpis":["efficiency","quality_score","delivery_time"]}'
redis-cli -p 6379 NODE.CREATE business_flow legal-dept department '{"name":"Legal Department","head":"Robert Wilson","employees":6,"budget":"$900K","location":"Floor 4","kpis":["compliance_score","contract_turnaround","risk_mitigation"]}'

echo "Creating business processes..."
redis-cli -p 6379 NODE.CREATE business_flow lead-qualification process '{"name":"Lead Qualification Process","owner":"sales-dept","duration":"2-5 days","automation_level":"70%","tools":["CRM","Lead Scoring","Email Automation"],"success_rate":"85%"}'
redis-cli -p 6379 NODE.CREATE business_flow proposal-creation process '{"name":"Proposal Creation","owner":"sales-dept","duration":"3-7 days","automation_level":"40%","tools":["Proposal Software","Template Library","E-signature"],"success_rate":"92%"}'
redis-cli -p 6379 NODE.CREATE business_flow contract-negotiation process '{"name":"Contract Negotiation","owner":"sales-dept","duration":"5-15 days","automation_level":"20%","tools":["Contract Management","Legal Review","Approval Workflow"],"success_rate":"78%"}'
redis-cli -p 6379 NODE.CREATE business_flow campaign-planning process '{"name":"Campaign Planning","owner":"marketing-dept","duration":"10-20 days","automation_level":"60%","tools":["Marketing Automation","Analytics","A/B Testing"],"success_rate":"88%"}'
redis-cli -p 6379 NODE.CREATE business_flow content-creation process '{"name":"Content Creation","owner":"marketing-dept","duration":"5-10 days","automation_level":"30%","tools":["CMS","Design Tools","Approval Workflow"],"success_rate":"95%"}'
redis-cli -p 6379 NODE.CREATE business_flow budget-approval process '{"name":"Budget Approval","owner":"finance-dept","duration":"3-10 days","automation_level":"80%","tools":["ERP","Approval Workflow","Financial Analytics"],"success_rate":"96%"}'
redis-cli -p 6379 NODE.CREATE business_flow invoice-processing process '{"name":"Invoice Processing","owner":"finance-dept","duration":"1-3 days","automation_level":"90%","tools":["OCR","AP Automation","ERP Integration"],"success_rate":"99%"}'
redis-cli -p 6379 NODE.CREATE business_flow employee-onboarding process '{"name":"Employee Onboarding","owner":"hr-dept","duration":"5-10 days","automation_level":"75%","tools":["HRIS","Digital Forms","Training Platform"],"success_rate":"93%"}'
redis-cli -p 6379 NODE.CREATE business_flow performance-review process '{"name":"Performance Review","owner":"hr-dept","duration":"15-30 days","automation_level":"50%","tools":["Performance Management","360 Feedback","Goal Tracking"],"success_rate":"87%"}'

echo "Creating decision points..."
redis-cli -p 6379 NODE.CREATE business_flow lead-scoring-gate decision '{"name":"Lead Scoring Gate","criteria":"score >= 75","automated":true,"fallback":"manual_review","approval_rate":"68%","avg_processing_time":"2h"}'
redis-cli -p 6379 NODE.CREATE business_flow budget-gate decision '{"name":"Budget Approval Gate","criteria":"amount <= $50K","automated":false,"approvers":["finance_head","cfo"],"approval_rate":"82%","avg_processing_time":"24h"}'
redis-cli -p 6379 NODE.CREATE business_flow legal-review-gate decision '{"name":"Legal Review Gate","criteria":"contract_value >= $100K","automated":false,"approvers":["legal_counsel"],"approval_rate":"91%","avg_processing_time":"48h"}'
redis-cli -p 6379 NODE.CREATE business_flow quality-gate decision '{"name":"Quality Assurance Gate","criteria":"quality_score >= 90%","automated":true,"fallback":"quality_team_review","approval_rate":"94%","avg_processing_time":"4h"}'
redis-cli -p 6379 NODE.CREATE business_flow compliance-gate decision '{"name":"Compliance Check Gate","criteria":"regulatory_compliance == true","automated":true,"fallback":"compliance_officer","approval_rate":"97%","avg_processing_time":"1h"}'

echo "Creating stakeholder roles..."
redis-cli -p 6379 NODE.CREATE business_flow sales-manager stakeholder '{"name":"Sales Manager","department":"sales-dept","responsibilities":["team_management","quota_oversight","deal_approval"],"authority_level":"department","max_approval":"$250K"}'
redis-cli -p 6379 NODE.CREATE business_flow marketing-director stakeholder '{"name":"Marketing Director","department":"marketing-dept","responsibilities":["strategy","budget_management","campaign_approval"],"authority_level":"department","max_approval":"$100K"}'
redis-cli -p 6379 NODE.CREATE business_flow cfo stakeholder '{"name":"Chief Financial Officer","department":"finance-dept","responsibilities":["financial_strategy","budget_approval","risk_management"],"authority_level":"executive","max_approval":"$1M"}'
redis-cli -p 6379 NODE.CREATE business_flow legal-counsel stakeholder '{"name":"Legal Counsel","department":"legal-dept","responsibilities":["contract_review","compliance","risk_assessment"],"authority_level":"specialist","max_approval":"legal_matters"}'
redis-cli -p 6379 NODE.CREATE business_flow hr-director stakeholder '{"name":"HR Director","department":"hr-dept","responsibilities":["talent_management","policy_development","employee_relations"],"authority_level":"department","max_approval":"$75K"}'

echo "Creating business documents..."
redis-cli -p 6379 NODE.CREATE business_flow sales-proposal document '{"name":"Sales Proposal","template":"proposal_v2.1","required_sections":["executive_summary","solution","pricing","timeline"],"approval_required":true,"retention":"7_years"}'
redis-cli -p 6379 NODE.CREATE business_flow service-contract document '{"name":"Service Contract","template":"service_agreement_v3.0","required_sections":["scope","deliverables","terms","payment"],"legal_review":true,"retention":"10_years"}'
redis-cli -p 6379 NODE.CREATE business_flow purchase-order document '{"name":"Purchase Order","template":"po_template_v1.5","required_fields":["vendor","items","amounts","delivery_date"],"approval_workflow":true,"retention":"5_years"}'
redis-cli -p 6379 NODE.CREATE business_flow marketing-brief document '{"name":"Marketing Brief","template":"campaign_brief_v2.0","required_sections":["objectives","target_audience","budget","timeline"],"stakeholder_review":true,"retention":"3_years"}'
redis-cli -p 6379 NODE.CREATE business_flow employee-handbook document '{"name":"Employee Handbook","template":"handbook_v4.2","required_sections":["policies","procedures","benefits","code_of_conduct"],"annual_review":true,"retention":"permanent"}'

echo "Creating systems and tools..."
redis-cli -p 6379 NODE.CREATE business_flow crm-system system '{"name":"Salesforce CRM","type":"customer_relationship_management","users":45,"license_cost":"$150/user/month","integration_apis":["marketing","finance","support"],"uptime_sla":"99.9%"}'
redis-cli -p 6379 NODE.CREATE business_flow erp-system system '{"name":"SAP ERP","type":"enterprise_resource_planning","users":80,"license_cost":"$200/user/month","modules":["finance","hr","procurement","inventory"],"uptime_sla":"99.95%"}'
redis-cli -p 6379 NODE.CREATE business_flow marketing-automation system '{"name":"HubSpot Marketing","type":"marketing_automation","users":20,"license_cost":"$800/month","features":["email_campaigns","lead_scoring","analytics"],"uptime_sla":"99.5%"}'
redis-cli -p 6379 NODE.CREATE business_flow document-management system '{"name":"SharePoint","type":"document_management","users":120,"license_cost":"$10/user/month","storage":"5TB","compliance_features":["version_control","audit_trail","retention_policies"]}'
redis-cli -p 6379 NODE.CREATE business_flow approval-workflow system '{"name":"Nintex Workflow","type":"business_process_automation","workflows":25,"license_cost":"$15000/year","features":["visual_designer","mobile_approval","analytics"],"integration_count":15}'

echo "Creating process flow connections - Lead to Sale..."
redis-cli -p 6379 EDGE.CREATE business_flow marketing-to-lead marketing-dept lead-qualification generates '{"lead_source":"campaigns","qualification_criteria":"BANT","handoff_sla":"24h","data_quality_score":"85%"}'
redis-cli -p 6379 EDGE.CREATE business_flow lead-to-scoring lead-qualification lead-scoring-gate evaluates '{"scoring_model":"predictive","factors":["company_size","budget","timeline","authority"],"threshold":"75/100"}'
redis-cli -p 6379 EDGE.CREATE business_flow scoring-to-proposal lead-scoring-gate proposal-creation triggers '{"condition":"qualified","assignment":"round_robin","sla":"48h","template_selection":"automated"}'
redis-cli -p 6379 EDGE.CREATE business_flow proposal-to-budget proposal-creation budget-gate requires '{"threshold":"$50K","auto_approval":"under_threshold","escalation":"finance_head","documentation_required":true}'
redis-cli -p 6379 EDGE.CREATE business_flow budget-to-legal budget-gate legal-review-gate escalates '{"contract_value":"$100K+","risk_assessment":"required","compliance_check":"automated","turnaround_sla":"48h"}'
redis-cli -p 6379 EDGE.CREATE business_flow legal-to-contract legal-review-gate contract-negotiation approves '{"terms_approved":"standard","custom_clauses":"reviewed","signature_authority":"delegated"}'

echo "Creating departmental dependencies..."
redis-cli -p 6379 EDGE.CREATE business_flow sales-to-marketing sales-dept marketing-dept collaborates '{"shared_goals":["lead_quality","pipeline_velocity"],"meetings":"weekly","shared_systems":["CRM","Marketing_Automation"],"kpi_alignment":"revenue_growth"}'
redis-cli -p 6379 EDGE.CREATE business_flow sales-to-finance sales-dept finance-dept reports '{"frequency":"monthly","metrics":["revenue","forecast","commission"],"approval_required":["discounts","payment_terms"],"system_integration":"CRM_to_ERP"}'
redis-cli -p 6379 EDGE.CREATE business_flow marketing-to-finance marketing-dept finance-dept requests '{"budget_approval":"quarterly","roi_reporting":"monthly","cost_allocation":"campaign_based","approval_threshold":"$25K"}'
redis-cli -p 6379 EDGE.CREATE business_flow hr-to-all hr-dept operations-dept supports '{"services":["recruitment","training","policy_enforcement"],"sla":"varies_by_service","escalation_path":"hr_director","compliance_oversight":"continuous"}'
redis-cli -p 6379 EDGE.CREATE business_flow legal-to-all legal-dept sales-dept advises '{"contract_review":"required","compliance_guidance":"ongoing","risk_assessment":"quarterly","approval_authority":"contracts_over_100K"}'

echo "Creating system integrations..."
redis-cli -p 6379 EDGE.CREATE business_flow crm-to-erp crm-system erp-system syncs '{"data_flow":"bidirectional","sync_frequency":"real_time","entities":["customers","orders","invoices"],"error_handling":"retry_with_alert"}'
redis-cli -p 6379 EDGE.CREATE business_flow marketing-to-crm marketing-automation crm-system integrates '{"lead_sync":"automatic","scoring_sync":"real_time","campaign_attribution":"tracked","data_quality":"validated"}'
redis-cli -p 6379 EDGE.CREATE business_flow workflow-to-systems approval-workflow crm-system orchestrates '{"approval_triggers":"automatic","status_updates":"real_time","notification_channels":["email","mobile","dashboard"],"audit_trail":"complete"}'

echo "Creating stakeholder involvement..."
redis-cli -p 6379 EDGE.CREATE business_flow manager-approves sales-manager budget-gate approves '{"authority_limit":"$250K","delegation_allowed":true,"escalation_required":"above_limit","approval_sla":"24h"}'
redis-cli -p 6379 EDGE.CREATE business_flow cfo-approves cfo budget-gate approves '{"authority_limit":"$1M","board_approval":"above_limit","financial_impact_assessment":"required","approval_sla":"48h"}'
redis-cli -p 6379 EDGE.CREATE business_flow legal-reviews legal-counsel legal-review-gate reviews '{"contract_complexity":"high","risk_assessment":"required","compliance_verification":"mandatory","review_sla":"48h"}'
redis-cli -p 6379 EDGE.CREATE business_flow hr-onboards hr-director employee-onboarding manages '{"process_ownership":"complete","system_access":"provisioned","compliance_training":"mandatory","completion_tracking":"automated"}'

echo "Creating document workflows..."
redis-cli -p 6379 EDGE.CREATE business_flow proposal-generates proposal-creation sales-proposal creates '{"template_selection":"automated","content_population":"crm_data","approval_routing":"hierarchical","version_control":"enabled"}'
redis-cli -p 6379 EDGE.CREATE business_flow contract-requires contract-negotiation service-contract generates '{"legal_template":"required","custom_terms":"negotiated","signature_workflow":"electronic","retention_compliance":"automated"}'
redis-cli -p 6379 EDGE.CREATE business_flow budget-creates budget-approval purchase-order generates '{"approval_evidence":"attached","vendor_verification":"required","delivery_tracking":"enabled","payment_terms":"net_30"}'

echo "✅ Business process flow graph created successfully!"
echo "📊 Graph contains:"
echo "   - 6 Organizational departments with KPIs and budgets"
echo "   - 9 Business processes with automation levels and success rates"
echo "   - 5 Decision gates with approval criteria and processing times"
echo "   - 5 Key stakeholder roles with authority levels"
echo "   - 5 Business document templates with retention policies"
echo "   - 5 Enterprise systems with integration capabilities"
echo "   - 25+ edges showing complex business relationships"
echo ""
echo "🔄 Business Flow Examples Created:"
echo "   1. Lead-to-Sale: Marketing → Lead Qualification → Scoring → Proposal → Budget → Legal → Contract"
echo "   2. Cross-Department: Sales ↔ Marketing ↔ Finance ↔ HR ↔ Legal (collaboration & dependencies)"
echo "   3. System Integration: CRM ↔ ERP ↔ Marketing Automation ↔ Workflow Engine"
echo "   4. Approval Hierarchy: Manager → Director → CFO → Board (authority levels)"
echo "   5. Document Lifecycle: Creation → Review → Approval → Signature → Retention"
