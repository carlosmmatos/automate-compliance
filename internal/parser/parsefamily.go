package parser

// parseFamily takes a string and returns an OpenControl friendly format for NIST Family.
func ParseFamily(family string) string {
	switch family {
	case "ACCESS_CONTROL":
		return "AC-Access_Control"
	case "AUDIT_AND_ACCOUNTABILITY":
		return "AU-Audit_and_Accountability"
	case "AWARENESS_AND_TRAINING":
		return "AT-Awareness_and_Training"
	case "CONFIGURATION_MANAGEMENT":
		return "CM-Configuration_Management"
	case "CONTINGENCY_PLANNING":
		return "CP-Contingency_Planning"
	case "IDENTIFICATION_AND_AUTHENTICATION":
		return "IA-Identification_and_Authentication"
	case "INCIDENT_RESPONSE":
		return "IR-Incident_Response"
	case "MAINTENANCE":
		return "MA-Maintenance"
	case "MEDIA_PROTECTION":
		return "MP-Media_Protection"
	case "PERSONNEL_SECURITY":
		return "PS-Personnel_Security"
	case "PHYSICAL_AND_ENVIRONMENTAL PROTECTION":
		return "PE-Physical_and_Environmental_Protection"
	case "PLANNING":
		return "PL-Planning"
	case "PROGRAM_MANAGEMENT":
		return "PM-Program_Management"
	case "RISK_ASSESSMENT":
		return "RA-Risk_Assessment"
	case "SECURITY_ASSESSMENT_AND_AUTHORIZATION":
		return "CA-Security_Assessment_and_Authorization"
	case "SYSTEM_AND_COMMUNICATIONS_PROTECTION":
		return "SC-System_and_Communications_Protection"
	case "SYSTEM_AND_INFORMATION_INTEGRITY":
		return "SI-System_and_Information_Integrity"
	case "SYSTEM_AND_SERVICES_ACQUISITION":
		return "SA-System_and_Services_Acquisition"
	default:
		return ""
	}
}
