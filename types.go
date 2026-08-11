package wardenauth

import "encoding/json"

type Effect string

const (
	EffectAllow Effect = "allow"
	EffectDeny  Effect = "deny"
)

type ScopeType string

const (
	ScopeTypeOrganization ScopeType = "organization"
	ScopeTypeWorkspace    ScopeType = "workspace"
	ScopeTypeApplication  ScopeType = "application"
)

type APIKeyType string

const (
	APIKeyTypeManagement  APIKeyType = "management"
	APIKeyTypeApplication APIKeyType = "application"
)

type BillingTier string

const (
	BillingTierFree       BillingTier = "free"
	BillingTierStarter    BillingTier = "starter"
	BillingTierGrowth     BillingTier = "growth"
	BillingTierBusiness   BillingTier = "business"
	BillingTierScale      BillingTier = "scale"
	BillingTierEnterprise BillingTier = "enterprise"
)

type PurchasableTier string

const (
	PurchasableTierStarter  PurchasableTier = "starter"
	PurchasableTierGrowth   PurchasableTier = "growth"
	PurchasableTierBusiness PurchasableTier = "business"
	PurchasableTierScale    PurchasableTier = "scale"
)

type SSOProtocol string

const (
	SSOProtocolSAML SSOProtocol = "saml"
	SSOProtocolOIDC SSOProtocol = "oidc"
)

type AuditExportFormat string

const (
	AuditExportCSV  AuditExportFormat = "csv"
	AuditExportJSON AuditExportFormat = "json"
)

type PaginatedResult[T any] struct {
	Items     []T    `json:"items"`
	NextToken string `json:"nextToken,omitempty"`
}

// ─── ABAC Conditions ──────────────────────────────────────────────────────

type AbacScalar any

type AbacConditionOperator string

type AbacConditionKind string

const (
	ConditionKindPredicate AbacConditionKind = "predicate"
	ConditionKindAll       AbacConditionKind = "all"
	ConditionKindAny       AbacConditionKind = "any"
	ConditionKindNot       AbacConditionKind = "not"
)

type AbacPredicateCondition struct {
	Kind     AbacConditionKind     `json:"kind"`
	Field    string                `json:"field"`
	Operator AbacConditionOperator `json:"operator"`
	Value    AbacScalar            `json:"value,omitempty"`
}

type AbacAllCondition struct {
	Kind       AbacConditionKind `json:"kind"`
	Conditions []AbacCondition   `json:"conditions"`
}

type AbacAnyCondition struct {
	Kind       AbacConditionKind `json:"kind"`
	Conditions []AbacCondition   `json:"conditions"`
}

type AbacNotCondition struct {
	Kind       AbacConditionKind `json:"kind"`
	Conditions []AbacCondition   `json:"conditions"`
}

type AbacCondition json.RawMessage

type AccessCheckContext map[string]any
type ScopeItem struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	ParentId *string `json:"parentId,omitempty"`
	Description *string `json:"description,omitempty"`
	AllowRoleInheritance *bool `json:"allowRoleInheritance,omitempty"`
	AccessScopeId *string `json:"accessScopeId,omitempty"`
}

type CreateScopeInput struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	ParentId *string `json:"parentId,omitempty"`
	Description *string `json:"description,omitempty"`
	InheritParent *bool `json:"inheritParent,omitempty"`
	AllowRoleInheritance *bool `json:"allowRoleInheritance,omitempty"`
	AccessScopeId *string `json:"accessScopeId,omitempty"`
}

type UpdateScopeInput struct {
	ID string `json:"id"`
	Name *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Type *string `json:"type,omitempty"`
	ParentId *string `json:"parentId,omitempty"`
	AllowRoleInheritance *bool `json:"allowRoleInheritance,omitempty"`
}

type PermissionItem struct {
	ID string `json:"id"`
	ScopeId string `json:"scopeId"`
	Resource string `json:"resource"`
	Action string `json:"action"`
	Effect string `json:"effect"`
	Name string `json:"name"`
	Description *string `json:"description,omitempty"`
	Conditions *interface{} `json:"conditions,omitempty"`
}

type CreatePermissionInput struct {
	ID string `json:"id"`
	ScopeId string `json:"scopeId"`
	Resource string `json:"resource"`
	Action string `json:"action"`
	Effect string `json:"effect"`
	Name string `json:"name"`
	Description *string `json:"description,omitempty"`
	Conditions *interface{} `json:"conditions,omitempty"`
}

type UpdatePermissionInput struct {
	Resource *string `json:"resource,omitempty"`
	Action *string `json:"action,omitempty"`
	Effect *string `json:"effect,omitempty"`
	Name *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Conditions *interface{} `json:"conditions,omitempty"`
}

type RoleItem struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description *string `json:"description,omitempty"`
	Scope *RoleItem `json:"scope"`
	Permissions []*RoleItem `json:"permissions"`
	InheritedFromScope *string `json:"inheritedFromScope,omitempty"`
}

type CreateRoleInput struct {
	ID string `json:"id"`
	ScopeId string `json:"scopeId"`
	Name string `json:"name"`
	Description *string `json:"description,omitempty"`
	PermissionIds []string `json:"permissionIds"`
	ParentRoleId *string `json:"parentRoleId,omitempty"`
}

type UpdateRoleInput struct {
	Name string `json:"name"`
	Description *string `json:"description,omitempty"`
	PermissionIds []string `json:"permissionIds"`
	ParentRoleId *string `json:"parentRoleId,omitempty"`
}

type AccessPolicyItem struct {
	ID string `json:"id"`
	SubjectId string `json:"subjectId"`
	Scope *AccessPolicyItem `json:"scope"`
	Roles []*AccessPolicyItem `json:"roles"`
	ExpiresAt *string `json:"expiresAt,omitempty"`
	ValidFrom *string `json:"validFrom,omitempty"`
}

type CreateAccessPolicyInput struct {
	SubjectId string `json:"subjectId"`
	ScopeId string `json:"scopeId"`
	RoleIds *[]string `json:"roleIds,omitempty"`
	PermissionIds *[]string `json:"permissionIds,omitempty"`
	SubjectType *string `json:"subjectType,omitempty"`
	ExpiresAt *string `json:"expiresAt,omitempty"`
	ValidFrom *string `json:"validFrom,omitempty"`
}

type UpdateAccessPolicyInput struct {
	RoleIds *[]string `json:"roleIds,omitempty"`
	ExpiresAt *string `json:"expiresAt,omitempty"`
	ValidFrom *string `json:"validFrom,omitempty"`
}

type ApiKeyItem struct {
	KeyId string `json:"keyId"`
	Type string `json:"type"`
	MaskedKey string `json:"maskedKey"`
	CreatedAt string `json:"createdAt"`
	SubjectId string `json:"subjectId"`
	Name string `json:"name"`
	ScopeId *string `json:"scopeId,omitempty"`
	KeyPrefix *string `json:"keyPrefix,omitempty"`
	ExpiresAt *string `json:"expiresAt,omitempty"`
	LastUsedAt *string `json:"lastUsedAt,omitempty"`
	Attributes *map[string]string `json:"attributes,omitempty"`
	AutoRotationEnabled *bool `json:"autoRotationEnabled,omitempty"`
	AutoRotationIntervalDays *float64 `json:"autoRotationIntervalDays,omitempty"`
	AutoRotationOverlapDays *float64 `json:"autoRotationOverlapDays,omitempty"`
	NextAutoRotationAt *string `json:"nextAutoRotationAt,omitempty"`
}

type ApiKeyCreatedItem struct {
	KeyId string `json:"keyId"`
	Type string `json:"type"`
	MaskedKey string `json:"maskedKey"`
	CreatedAt string `json:"createdAt"`
	SubjectId string `json:"subjectId"`
	Name string `json:"name"`
	ScopeId *string `json:"scopeId,omitempty"`
	KeyPrefix *string `json:"keyPrefix,omitempty"`
	ExpiresAt *string `json:"expiresAt,omitempty"`
	LastUsedAt *string `json:"lastUsedAt,omitempty"`
	Attributes *map[string]string `json:"attributes,omitempty"`
	AutoRotationEnabled *bool `json:"autoRotationEnabled,omitempty"`
	AutoRotationIntervalDays *float64 `json:"autoRotationIntervalDays,omitempty"`
	AutoRotationOverlapDays *float64 `json:"autoRotationOverlapDays,omitempty"`
	NextAutoRotationAt *string `json:"nextAutoRotationAt,omitempty"`
	RawKey string `json:"rawKey"`
}

type CreateApiKeyInput struct {
	Name string `json:"name"`
	Type *string `json:"type,omitempty"`
	KeyPrefix *string `json:"keyPrefix,omitempty"`
	ExpiresAt *string `json:"expiresAt,omitempty"`
	Attributes *map[string]string `json:"attributes,omitempty"`
	Permissions *[]string `json:"permissions,omitempty"`
	Roles *[]string `json:"roles,omitempty"`
	RateLimit *int `json:"rateLimit,omitempty"`
	AutoRotationEnabled *bool `json:"autoRotationEnabled,omitempty"`
	AutoRotationIntervalDays *int `json:"autoRotationIntervalDays,omitempty"`
	AutoRotationOverlapDays *int `json:"autoRotationOverlapDays,omitempty"`
}

type ApiKeyRotationResult struct {
	KeyId string `json:"keyId"`
	MaskedKey string `json:"maskedKey"`
	RawKey string `json:"rawKey"`
	Name string `json:"name"`
	OldKeyId string `json:"oldKeyId"`
	OverlapExpiresAt *string `json:"overlapExpiresAt,omitempty"`
}

type ApiKeyRevealRotationResult struct {
	KeyId string `json:"keyId"`
	ApiKey string `json:"apiKey"`
}

type RotateApiKeyInput struct {
	OverlapDays *int `json:"overlapDays,omitempty"`
}

type UpdateApiKeyAutoRotationInput struct {
	Enabled bool `json:"enabled"`
	IntervalDays *float64 `json:"intervalDays,omitempty"`
	OverlapDays *float64 `json:"overlapDays,omitempty"`
}

type WebhookEndpointItem struct {
	ID string `json:"id"`
	ScopeId string `json:"scopeId"`
	Url string `json:"url"`
	Events []string `json:"events"`
	Active bool `json:"active"`
	Name *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt,omitempty"`
	DisabledAt *string `json:"disabledAt,omitempty"`
	LastNotifiedLevel *string `json:"lastNotifiedLevel,omitempty"`
}

type WebhookEndpointCreatedItem struct {
	ID string `json:"id"`
	ScopeId string `json:"scopeId"`
	Url string `json:"url"`
	Events []string `json:"events"`
	Active bool `json:"active"`
	Name *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt,omitempty"`
	DisabledAt *string `json:"disabledAt,omitempty"`
	LastNotifiedLevel *string `json:"lastNotifiedLevel,omitempty"`
	Secret string `json:"secret"`
}

type CreateWebhookEndpointInput struct {
	ID *string `json:"id,omitempty"`
	Url string `json:"url"`
	Events []string `json:"events"`
	Active *bool `json:"active,omitempty"`
	Name *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateWebhookEndpointInput struct {
	Url *string `json:"url,omitempty"`
	Events *[]string `json:"events,omitempty"`
	Active *bool `json:"active,omitempty"`
	Name *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type WebhookRotateSecretResult struct {
	Secret string `json:"secret"`
}

type WebhookDeliveryItem struct {
	ID string `json:"id"`
	EndpointId string `json:"endpointId"`
	ScopeId string `json:"scopeId"`
	Url string `json:"url"`
	EventType string `json:"eventType"`
	Success bool `json:"success"`
	HttpStatus *float64 `json:"httpStatus,omitempty"`
	DurationMs float64 `json:"durationMs"`
	Attempts float64 `json:"attempts"`
	ErrorMessage *string `json:"errorMessage,omitempty"`
	RequestBody *string `json:"requestBody,omitempty"`
	ResponseBody *string `json:"responseBody,omitempty"`
	DeliveredAt string `json:"deliveredAt"`
	SkipReason *string `json:"skipReason,omitempty"`
	NextRetryAt *string `json:"nextRetryAt,omitempty"`
	ReceiveCount *float64 `json:"receiveCount,omitempty"`
	PayloadData *string `json:"payloadData,omitempty"`
	Ttl *float64 `json:"ttl,omitempty"`
}

type AccessCheckInput struct {
	SubjectId     string                  `json:"subjectId"`
	ScopeId       string                  `json:"scopeId"`
	Resource      string                  `json:"resource"`
	ResourceId    *string                 `json:"resourceId,omitempty"`
	Action        string                  `json:"action"`
	Context       *map[string]interface{} `json:"context,omitempty"`
	IncludeReason *bool                   `json:"includeReason,omitempty"`
}

type AccessCheckReasoningEntry struct {
	RoleId string `json:"roleId"`
	RoleName string `json:"roleName"`
	Effect *string `json:"effect,omitempty"`
	PermissionId string `json:"permissionId"`
	ScopeId string `json:"scopeId"`
	InheritedFromScope *string `json:"inheritedFromScope,omitempty"`
}

type AccessCheckReasoning struct {
	MatchedBy []*AccessCheckReasoningEntry `json:"matchedBy"`
	DeniedBy []*AccessCheckReasoningEntry `json:"deniedBy"`
	ScopeChain []string `json:"scopeChain"`
}

type AccessCheckResult struct {
	SubjectId string `json:"subjectId"`
	ScopeId string `json:"scopeId"`
	Resource string `json:"resource"`
	Action string `json:"action"`
	Context *map[string]interface{} `json:"context,omitempty"`
	IncludeReason *bool `json:"includeReason,omitempty"`
	Allowed bool `json:"allowed"`
	Error *string `json:"error,omitempty"`
	InheritedFromScope *string `json:"inheritedFromScope,omitempty"`
	Reasoning *AccessCheckReasoning `json:"reasoning,omitempty"`
}

type ReceiptIssueInput struct {
	SubjectId string `json:"subjectId"`
	Resource string `json:"resource"`
	Action string `json:"action"`
}

type ReceiptIssueResult struct {
	Receipt string `json:"receipt"`
	Decision AccessCheckResult `json:"decision"`
}

type ReceiptVerifyInput struct {
	Receipt string `json:"receipt"`
}

type ReceiptVerifyResult struct {
	Valid bool `json:"valid"`
	Claims *ReceiptClaims `json:"claims,omitempty"`
	Reason *string `json:"reason,omitempty"`
}

type ReceiptClaims struct {
	Jti *string `json:"jti,omitempty"`
	SubjectId *string `json:"subjectId,omitempty"`
	ScopeId *string `json:"scopeId,omitempty"`
	Resource *string `json:"resource,omitempty"`
	Action *string `json:"action,omitempty"`
	Allowed *bool `json:"allowed,omitempty"`
	IssuedAt *float64 `json:"issuedAt,omitempty"`
	ExpiresAt *float64 `json:"expiresAt,omitempty"`
}

type ListPermissionsInput struct {
	SubjectId string `json:"subjectId"`
	ScopeId string `json:"scopeId"`
}

type ListPermissionsResult struct {
	Permissions []*ListPermissionsResult `json:"permissions"`
}

type ListRolesInput struct {
	SubjectId string `json:"subjectId"`
	ScopeId string `json:"scopeId"`
}

type ListRolesResult struct {
	Roles []*ListRolesResult `json:"roles"`
}

type SimulateAccessPermissionInput struct {
	ID string `json:"id"`
	Resource string `json:"resource"`
	Action string `json:"action"`
	Effect string `json:"effect"`
	Description *string `json:"description,omitempty"`
}

type SimulateAccessRoleInput struct {
	ID string `json:"id"`
	Name *string `json:"name,omitempty"`
	PermissionIds []string `json:"permissionIds"`
}

type SimulateAccessCheckInput struct {
	Resource string `json:"resource"`
	Action string `json:"action"`
	Context *map[string]interface{} `json:"context,omitempty"`
}

type SimulateAccessInput struct {
	SubjectId string `json:"subjectId"`
	ScopeId string `json:"scopeId"`
	Checks []*SimulateAccessInput `json:"checks"`
	Permissions []*SimulateAccessInput `json:"permissions"`
	Roles []*SimulateAccessInput `json:"roles"`
	SubjectRoleIds []string `json:"subjectRoleIds"`
	IncludeReason *bool `json:"includeReason,omitempty"`
}

type SimulateAccessResultEntry struct {
	Resource string `json:"resource"`
	Action string `json:"action"`
	Allowed bool `json:"allowed"`
	Reasoning **SimulateAccessResultEntry `json:"reasoning,omitempty"`
}

type SimulateAccessResult struct {
	SubjectId string `json:"subjectId"`
	ScopeId string `json:"scopeId"`
	Results []*SimulateAccessResult `json:"results"`
}

type AuditLogItem struct {
	ID string `json:"id"`
	EventType string `json:"eventType"`
	Timestamp string `json:"timestamp"`
	Actor string `json:"actor"`
	ScopeId string `json:"scopeId"`
	ResourceType *string `json:"resourceType,omitempty"`
	ResourceId *string `json:"resourceId,omitempty"`
	TargetSubjectId *string `json:"targetSubjectId,omitempty"`
	Action *string `json:"action,omitempty"`
	Result string `json:"result"`
	Metadata *map[string]interface{} `json:"metadata,omitempty"`
	IpAddress *string `json:"ipAddress,omitempty"`
	TraceId *string `json:"traceId,omitempty"`
	UserAgent *string `json:"userAgent,omitempty"`
}

type AuditExportInput struct {
	ScopeId *string `json:"scopeId,omitempty"`
	ActorId *string `json:"actorId,omitempty"`
	ResourceType *string `json:"resourceType,omitempty"`
	ResourceId *string `json:"resourceId,omitempty"`
	StartDate *string `json:"startDate,omitempty"`
	EndDate *string `json:"endDate,omitempty"`
	EventTypes *[]string `json:"eventTypes,omitempty"`
	MaxRows *float64 `json:"maxRows,omitempty"`
	Format *string `json:"format,omitempty"`
}

type AuditVerifyInput struct {
	ScopeId string `json:"scopeId"`
	StartDate *string `json:"startDate,omitempty"`
	EndDate *string `json:"endDate,omitempty"`
	MaxRecords *float64 `json:"maxRecords,omitempty"`
}

type AuditVerifyResult struct {
	Total float64 `json:"total"`
	Matched float64 `json:"matched"`
	MismatchCount float64 `json:"mismatchCount"`
	Mismatches []*AuditVerifyResult `json:"mismatches"`
}

type ResourceTypeItem struct {
	ID string `json:"id"`
	ScopeId string `json:"scopeId"`
	Name string `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedAt string `json:"createdAt"`
}

type CreateResourceTypeInput struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description *string `json:"description,omitempty"`
}

type SubscriptionSummary struct {
	Tier string `json:"tier"`
	TierName string `json:"tierName"`
	MonthlyLimit *float64 `json:"monthlyLimit,omitempty"`
	CurrentMonthCount float64 `json:"currentMonthCount"`
	PercentUsed *float64 `json:"percentUsed,omitempty"`
	NotificationEmails []string `json:"notificationEmails"`
	AlertThresholdPercent float64 `json:"alertThresholdPercent"`
	AnomalyAlertsEnabled bool `json:"anomalyAlertsEnabled"`
	BillingInterval *string `json:"billingInterval,omitempty"`
	StripeStatus *string `json:"stripeStatus,omitempty"`
}

type UpdateSubscriptionInput struct {
	NotificationEmails *[]string `json:"notificationEmails,omitempty"`
	AlertThresholdPercent *float64 `json:"alertThresholdPercent,omitempty"`
	AnomalyAlertsEnabled *bool `json:"anomalyAlertsEnabled,omitempty"`
}

type CreateCheckoutSessionInput struct {
	Tier string `json:"tier"`
	Interval string `json:"interval"`
	CustomerEmail *string `json:"customerEmail,omitempty"`
	SuccessUrl string `json:"successUrl"`
	CancelUrl string `json:"cancelUrl"`
}

type CreateCheckoutSessionResult struct {
	Url string `json:"url"`
}

type CreateCustomerPortalSessionResult struct {
	Url string `json:"url"`
}

type OverageStatus struct {
	Tier string `json:"tier"`
	TierName string `json:"tierName"`
	IncludedChecks float64 `json:"includedChecks"`
	CurrentMonthChecks float64 `json:"currentMonthChecks"`
	OverageChecks float64 `json:"overageChecks"`
	OverageRatePerMillion float64 `json:"overageRatePerMillion"`
	EstimatedOverageCharge float64 `json:"estimatedOverageCharge"`
	IsPaused bool `json:"isPaused"`
	ConsentCount float64 `json:"consentCount"`
	NextPauseAtPercent float64 `json:"nextPauseAtPercent"`
	NextPauseAtChecks float64 `json:"nextPauseAtChecks"`
	IsHardPause bool `json:"isHardPause"`
	BillingInterval *string `json:"billingInterval,omitempty"`
	StripeStatus *string `json:"stripeStatus,omitempty"`
}

type GrantOverageConsentInput struct {
	IpAddress *string `json:"ipAddress,omitempty"`
}

type GrantOverageConsentResult struct {
	Resumed bool `json:"resumed"`
	Tier string `json:"tier"`
	ConsentCount float64 `json:"consentCount"`
	NextPauseAtPercent float64 `json:"nextPauseAtPercent"`
	NextPauseAtChecks float64 `json:"nextPauseAtChecks"`
}

type SodConstraintItem struct {
	ID string `json:"id"`
	ScopeId string `json:"scopeId"`
	RoleIdA string `json:"roleIdA"`
	RoleIdB string `json:"roleIdB"`
	Description *string `json:"description,omitempty"`
}

type CreateSodConstraintInput struct {
	ID *string `json:"id,omitempty"`
	RoleIdA string `json:"roleIdA"`
	RoleIdB string `json:"roleIdB"`
	Description *string `json:"description,omitempty"`
}

type WardenAuthClientConfig struct {
	ApiUrl string `json:"apiUrl"`
	ApiKey string `json:"apiKey"`
}

type AuthorizationManifestPermissionInput struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Resource string `json:"resource"`
	Action string `json:"action"`
	Effect string `json:"effect"`
	Description *string `json:"description,omitempty"`
}

type AuthorizationManifestRoleInput struct {
	ID string `json:"id"`
	Name string `json:"name"`
	PermissionIds []string `json:"permissionIds"`
	Description *string `json:"description,omitempty"`
	IsTemplate *bool `json:"isTemplate,omitempty"`
}

type AuthorizationManifestAccessPolicyInput struct {
	SubjectId string `json:"subjectId"`
	SubjectType *string `json:"subjectType,omitempty"`
	RoleIds []string `json:"roleIds"`
	ExpiresAt *string `json:"expiresAt,omitempty"`
	ValidFrom *string `json:"validFrom,omitempty"`
}

type AuthorizationManifestSpecInput struct {
	Mode string `json:"mode"`
	DeletionPolicy string `json:"deletionPolicy"`
	DryRun *bool `json:"dryRun,omitempty"`
	IdempotencyKey *string `json:"idempotencyKey,omitempty"`
	RequestTimestamp *string `json:"requestTimestamp,omitempty"`
	Permissions []*AuthorizationManifestSpecInput `json:"permissions"`
	Roles []*AuthorizationManifestSpecInput `json:"roles"`
	AccessPolicies []*AuthorizationManifestSpecInput `json:"accessPolicies"`
}

type AuthorizationManifestApplyInput struct {
	ApiVersion string `json:"apiVersion"`
	Kind string `json:"kind"`
	Spec *AuthorizationManifestApplyInput `json:"spec"`
}

type ApplyScopeManifestInput struct {
	Manifest *ApplyScopeManifestInput `json:"manifest"`
	Serialization *string `json:"serialization,omitempty"`
}

type ApplyScopeManifestResultSummary struct {
	TotalPlanned float64 `json:"totalPlanned"`
	Applied float64 `json:"applied"`
	Failed float64 `json:"failed"`
	Planned float64 `json:"planned"`
}

type ApplyScopeManifestOperationResult struct {
	ResourceType string `json:"resourceType"`
	Operation string `json:"operation"`
	ResourceKey string `json:"resourceKey"`
	Status string `json:"status"`
	Error *string `json:"error,omitempty"`
}

type ApplyScopeManifestResult struct {
	ScopeId string `json:"scopeId"`
	DryRun bool `json:"dryRun"`
	IdempotencyKey *string `json:"idempotencyKey,omitempty"`
	ManifestHash string `json:"manifestHash"`
	Summary *ApplyScopeManifestResult `json:"summary"`
	Operations []*ApplyScopeManifestResult `json:"operations"`
}

type TeamMemberItem struct {
	SubjectId string `json:"subjectId"`
	ScopeId string `json:"scopeId"`
	RoleIds []string `json:"roleIds"`
	Email *string `json:"email,omitempty"`
	Name *string `json:"name,omitempty"`
	InvitedBy string `json:"invitedBy"`
	CreatedAt string `json:"createdAt"`
	Status string `json:"status"`
}

type AddTeamMemberInput struct {
	SubjectId *string `json:"subjectId,omitempty"`
	RoleIds []string `json:"roleIds"`
	Email *string `json:"email,omitempty"`
	Name *string `json:"name,omitempty"`
}

type MintSessionTokenInput struct {
	TtlSeconds *int `json:"ttlSeconds,omitempty"`
	PermissionIds *[]string `json:"permissionIds,omitempty"`
	Purpose *string `json:"purpose,omitempty"`
}

type MintSessionTokenResult struct {
	Token string `json:"token"`
	TokenType string `json:"tokenType"`
	PrincipalId string `json:"principalId"`
	ScopeId string `json:"scopeId"`
	Jti string `json:"jti"`
	IssuedAt float64 `json:"issuedAt"`
	ExpiresAt string `json:"expiresAt"`
	PermissionIds *[]string `json:"permissionIds,omitempty"`
	Purpose *string `json:"purpose,omitempty"`
}

type AnomalyResult struct {
	ScopeId string `json:"scopeId"`
	Metric string `json:"metric"`
	CurrentValue float64 `json:"currentValue"`
	AverageValue float64 `json:"averageValue"`
	Threshold float64 `json:"threshold"`
	Flagged bool `json:"flagged"`
	Period string `json:"period"`
	ID *string `json:"id,omitempty"`
	Proof *[]*AnomalyResult `json:"proof,omitempty"`
}

type DashboardStats struct {
	ScopeName string `json:"scopeName"`
	Counts *DashboardStats `json:"counts"`
	Usage *DashboardStats `json:"usage"`
	Organization **DashboardStats `json:"organization,omitempty"`
	AnomalyCounts **DashboardStats `json:"anomalyCounts,omitempty"`
	Anomalies []*DashboardStats `json:"anomalies"`
	RecentActivity []*DashboardStats `json:"recentActivity"`
}

type IpAllowlistResult struct {
	Cidrs []string `json:"cidrs"`
}

type UpdateIpAllowlistInput struct {
	Cidrs []string `json:"cidrs"`
}

type TierActions struct {
	Low string `json:"low"`
	Medium string `json:"medium"`
	High string `json:"high"`
}

type TierPolicyResult struct {
	Tiers TierActions `json:"tiers"`
}

type UpdateTierPolicyInput struct {
	Low string `json:"low"`
	Medium string `json:"medium"`
	High string `json:"high"`
}

type OrganizationItem struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Description *string `json:"description,omitempty"`
	ParentId *string `json:"parentId,omitempty"`
	Billing **OrganizationItem `json:"billing,omitempty"`
}

type UpdateOrganizationInput struct {
	Name *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type SsoConfigItem struct {
	OrgId string `json:"orgId"`
	Provider string `json:"provider"`
	Protocol string `json:"protocol"`
	CognitoIdpName string `json:"cognitoIdpName"`
	Status string `json:"status"`
	SsoLoginUrl string `json:"ssoLoginUrl"`
	ClientId *string `json:"clientId,omitempty"`
	IssuerUrl *string `json:"issuerUrl,omitempty"`
	SpMetadataUrl *string `json:"spMetadataUrl,omitempty"`
	CallbackUrl *string `json:"callbackUrl,omitempty"`
}

type CreateSsoConfigInput struct {
	Provider string `json:"provider"`
	Protocol string `json:"protocol"`
	MetadataUrl *string `json:"metadataUrl,omitempty"`
	MetadataXml *string `json:"metadataXml,omitempty"`
	ClientId *string `json:"clientId,omitempty"`
	ClientSecret *string `json:"clientSecret,omitempty"`
	IssuerUrl *string `json:"issuerUrl,omitempty"`
}

type UpdateSsoConfigInput struct {
	MetadataUrl *string `json:"metadataUrl,omitempty"`
	MetadataXml *string `json:"metadataXml,omitempty"`
	ClientId *string `json:"clientId,omitempty"`
	ClientSecret *string `json:"clientSecret,omitempty"`
	IssuerUrl *string `json:"issuerUrl,omitempty"`
}

type SsoTestResult struct {
	Success bool `json:"success"`
	Error *string `json:"error,omitempty"`
}

type RotateScimTokenResult struct {
	Token string `json:"token"`
	CreatedAt string `json:"createdAt"`
}

type ScimConfig struct {
	OrgId string `json:"orgId"`
	ScimBaseUrl string `json:"scimBaseUrl"`
	BearerToken string `json:"bearerToken"`
}

type UserIdentityItem struct {
	SubjectId      string  `json:"subjectId"`
	Email          *string `json:"email,omitempty"`
	ProvisionedVia *string `json:"provisionedVia,omitempty"`
	Status         *string `json:"status,omitempty"`
	ExternalId     *string `json:"externalId,omitempty"`
}

type WorkspaceUsageItem struct {
	ScopeId string `json:"scopeId"`
	ScopeName string `json:"scopeName"`
	ChecksToday float64 `json:"checksToday"`
	GrantedToday float64 `json:"grantedToday"`
	DeniedToday float64 `json:"deniedToday"`
	ChecksThisMonth float64 `json:"checksThisMonth"`
}

type WorkspaceUsage struct {
	Items []*WorkspaceUsage `json:"items"`
	TotalChecksToday float64 `json:"totalChecksToday"`
	TotalGrantedToday float64 `json:"totalGrantedToday"`
	TotalDeniedToday float64 `json:"totalDeniedToday"`
	TotalChecksThisMonth float64 `json:"totalChecksThisMonth"`
}

type ImportResult struct {
	Created float64  `json:"created"`
	Skipped float64  `json:"skipped"`
	Errors  []string `json:"errors"`
}

type TupleKey struct {
	Subject  string `json:"subject"`
	Relation string `json:"relation"`
	Object   string `json:"object"`
}

type TupleWriteEntry struct {
	Subject  string            `json:"subject"`
	Relation string            `json:"relation"`
	Object   string            `json:"object"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type TupleWriteInput struct {
	Writes  []TupleWriteEntry `json:"writes,omitempty"`
	Deletes []TupleKey        `json:"deletes,omitempty"`
}

type TupleWriteResult struct {
	Written float64 `json:"written"`
	Deleted float64 `json:"deleted"`
}

type TupleItem struct {
	Subject  string            `json:"subject"`
	Relation string            `json:"relation"`
	Object   string            `json:"object"`
	ScopeId  string            `json:"scopeId"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type TupleListResult struct {
	Tuples []TupleItem `json:"tuples"`
}

type AgentIntentInstance struct {
	ResourceType string `json:"resourceType"`
	ResourceId   string `json:"resourceId"`
}

type AgentIntent struct {
	IntegrationId string                `json:"integrationId"`
	Tools         []string              `json:"tools"`
	Instances     []AgentIntentInstance `json:"instances,omitempty"`
}

type MintIntentSessionTokenInput struct {
	TtlSeconds *int        `json:"ttlSeconds,omitempty"`
	Purpose    *string     `json:"purpose,omitempty"`
	Intent     AgentIntent `json:"intent"`
}

type MintIntentSessionTokenResult struct {
	Token *string `json:"token,omitempty"`
	TokenType *string `json:"tokenType,omitempty"`
	PrincipalId *string `json:"principalId,omitempty"`
	ScopeId *string `json:"scopeId,omitempty"`
	Jti *string `json:"jti,omitempty"`
	IssuedAt *float64 `json:"issuedAt,omitempty"`
	ExpiresAt *string `json:"expiresAt,omitempty"`
	IntentHash *string `json:"intentHash,omitempty"`
	Intent *AgentIntent `json:"intent,omitempty"`
	Purpose *string `json:"purpose,omitempty"`
	ApprovalRequired *bool `json:"approvalRequired,omitempty"`
	Approval *ApprovalItem `json:"approval,omitempty"`
	CriticalTools *[]string `json:"criticalTools,omitempty"`
}

type ApprovalItem struct {
	ID string `json:"id"`
	ScopeId string `json:"scopeId"`
	RequesterId string `json:"requesterId"`
	SubjectId string `json:"subjectId"`
	Kind string `json:"kind"`
	Action string `json:"action"`
	ProposedRoleIds []string `json:"proposedRoleIds"`
	PreviousRoleIds *[]string `json:"previousRoleIds,omitempty"`
	Status string `json:"status"`
	ApproverId *string `json:"approverId,omitempty"`
	Reason *string `json:"reason,omitempty"`
	DeniedReason *string `json:"deniedReason,omitempty"`
	ValidFrom *string `json:"validFrom,omitempty"`
	ExpiresAt *string `json:"expiresAt,omitempty"`
	CreatedAt string `json:"createdAt"`
	ResolvedAt *string `json:"resolvedAt,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type IntegrationStatus string

const (
	IntegrationStatusActive   IntegrationStatus = "active"
	IntegrationStatusDisabled IntegrationStatus = "disabled"
)

type McpToolConfig struct {
	Name string `json:"name"`
	Tier string `json:"tier"`
	// Hitl, if set, requires human approval: "self" (requester, audited) or "admin" (separation of duties).
	Hitl        *string `json:"hitl,omitempty"`
	Description *string `json:"description,omitempty"`
}

type Integration struct {
	ScopeId     string            `json:"scopeId"`
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Provider    string            `json:"provider"`
	Status      IntegrationStatus `json:"status"`
	UpstreamUrl *string           `json:"upstreamUrl,omitempty"`
	CreatedAt   float64           `json:"createdAt"`
	CreatedBy   *string           `json:"createdBy,omitempty"`
	Tools       []McpToolConfig   `json:"tools,omitempty"`
}

type CreateIntegrationInput struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Provider    string             `json:"provider"`
	Status      *IntegrationStatus `json:"status,omitempty"`
	UpstreamUrl *string            `json:"upstreamUrl,omitempty"`
	Tools       []McpToolConfig    `json:"tools,omitempty"`
}

type UpdateIntegrationInput struct {
	Name        *string            `json:"name,omitempty"`
	Status      *IntegrationStatus `json:"status,omitempty"`
	UpstreamUrl *string            `json:"upstreamUrl,omitempty"`
	Tools       []McpToolConfig    `json:"tools,omitempty"`
}

type VerifyIntentCallInstanceInput struct {
	ResourceType string `json:"resourceType"`
	ResourceId   string `json:"resourceId"`
}

type VerifyIntentCallInput struct {
	Token    string                         `json:"token"`
	Tool     string                         `json:"tool"`
	Instance *VerifyIntentCallInstanceInput `json:"instance,omitempty"`
}

type VerifyIntentCallResult struct {
	Allowed     bool    `json:"allowed"`
	Jti         string  `json:"jti"`
	PrincipalId string  `json:"principalId"`
	Reason      *string `json:"reason,omitempty"`
	Pending     *bool   `json:"pending,omitempty"`
	ApprovalId  *string `json:"approvalId,omitempty"`
}

type McpConsentServerSummary struct {
	ServerKey string     `json:"serverKey"`
	Name      string     `json:"name"`
	UpstreamURL *string  `json:"upstreamUrl,omitempty"`
	ToolCount int        `json:"toolCount"`
}

type McpConsentToolView struct {
	Name        string  `json:"name"`
	Tier        string  `json:"tier"`
	Description *string `json:"description,omitempty"`
	HITL        *string `json:"hitl,omitempty"`
	Allowed     bool    `json:"allowed"`
}

type McpConsentContext struct {
	Server       McpConsentServerSummary `json:"server"`
	MaxTier      *string                 `json:"maxTier,omitempty"`
	Tools        []McpConsentToolView    `json:"tools"`
	AllowedCount int                     `json:"allowedCount"`
	DeniedCount  int                     `json:"deniedCount"`
}

type McpConsentServerSelection struct {
	ServerKey string `json:"serverKey"`
	Tier      string `json:"tier"`
}

type McpConsentGrantBody struct {
	Servers         []McpConsentServerSelection `json:"servers"`
	DurationSeconds *int                        `json:"durationSeconds,omitempty"`
	AuthRequestID   *string                     `json:"authRequestId,omitempty"`
	ScopeID         *string                     `json:"scopeId,omitempty"`
}

type McpConsentGrantResult struct {
	Status      string  `json:"status"`
	Token       *string `json:"token,omitempty"`
	RedirectURL *string `json:"redirectUrl,omitempty"`
}

type McpGrantServerSummary struct {
	ServerKey string `json:"serverKey"`
	Tier      string `json:"tier"`
	ToolCount int    `json:"toolCount"`
}

type McpGrantSummary struct {
	GrantID   string                    `json:"grantId"`
	Servers   []McpGrantServerSummary   `json:"servers"`
	ClientID  *string                   `json:"clientId,omitempty"`
	CreatedAt float64                   `json:"createdAt"`
	ExpiresAt *float64                  `json:"expiresAt,omitempty"`
}

type McpCallContext struct {
	Action  string                       `json:"action"`
	Targets []McpCallContextTargetEntry  `json:"targets"`
}

type McpCallContextTargetEntry struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type McpApprovalSummary struct {
	ID        string          `json:"id"`
	SubjectID string          `json:"subjectId"`
	ServerKey string          `json:"serverKey"`
	Tool      string          `json:"tool"`
	Level     string          `json:"level"`
	Context   *McpCallContext `json:"context,omitempty"`
	CreatedAt float64         `json:"createdAt"`
	ExpiresAt *float64        `json:"expiresAt,omitempty"`
}

type McpApprovalHistoryItem struct {
	ID           string          `json:"id"`
	SubjectID    string          `json:"subjectId"`
	ServerKey    string          `json:"serverKey"`
	Tool         string          `json:"tool"`
	Level        string          `json:"level"`
	Context      *McpCallContext `json:"context,omitempty"`
	CreatedAt    float64         `json:"createdAt"`
	ExpiresAt    *float64        `json:"expiresAt,omitempty"`
	Status       string          `json:"status"`
	DecidedBy    *string         `json:"decidedBy,omitempty"`
	DecidedAt    *float64        `json:"decidedAt,omitempty"`
	SelfApproved *bool           `json:"selfApproved,omitempty"`
	Reason       *string         `json:"reason,omitempty"`
}

type McpVelocityEffective struct {
	Enabled            bool `json:"enabled"`
	PerGrantPerMinute  int  `json:"perGrantPerMinute"`
	PerServerPerMinute int  `json:"perServerPerMinute"`
	PerToolPerMinute   int  `json:"perToolPerMinute"`
}

type McpVelocityOverrides struct {
	Enabled            bool `json:"enabled"`
	PerGrantPerMinute  int  `json:"perGrantPerMinute"`
	PerServerPerMinute int  `json:"perServerPerMinute"`
	PerToolPerMinute   int  `json:"perToolPerMinute"`
}

type McpVelocityConfig struct {
	Effective McpVelocityEffective  `json:"effective"`
	Overrides *McpVelocityOverrides `json:"overrides,omitempty"`
}

type AgentIdentifyInput struct {
	DelegatingUserID string      `json:"delegatingUserId"`
	WorkflowID       *string     `json:"workflowId,omitempty"`
	Intent           AgentIntent `json:"intent"`
	TTLSeconds       *int        `json:"ttlSeconds,omitempty"`
	Purpose          *string     `json:"purpose,omitempty"`
}

type AgentIdentifyResult struct {
	Token                  string      `json:"token"`
	TokenType              string      `json:"tokenType"`
	PrincipalID            string      `json:"principalId"`
	ScopeID                string      `json:"scopeId"`
	JTI                    string      `json:"jti"`
	IssuedAt               int         `json:"issuedAt"`
	ExpiresAt              string      `json:"expiresAt"`
	IntentHash             string      `json:"intentHash"`
	Intent                 AgentIntent `json:"intent"`
	DelegatingUserID       string      `json:"delegatingUserId"`
	WorkflowID             *string     `json:"workflowId,omitempty"`
	EffectivePermissionIDs []string    `json:"effectivePermissionIds"`
	Purpose                *string     `json:"purpose,omitempty"`
}

type AgentCheckInput struct {
	IdentityToken string               `json:"identityToken"`
	Tool          string               `json:"tool"`
	Instance      *AgentIntentInstance `json:"instance,omitempty"`
}

type AgentCheckResult struct {
	Allowed       bool    `json:"allowed"`
	JTI           string  `json:"jti"`
	PrincipalID   string  `json:"principalId"`
	Reason        *string `json:"reason,omitempty"`
	IntegrationID *string `json:"integrationId,omitempty"`
	Pending       *bool   `json:"pending,omitempty"`
	ApprovalID    *string `json:"approvalId,omitempty"`
}

