-- +goose Up
CREATE TABLE organizations (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(120) NOT NULL,
    name VARCHAR(150) NOT NULL,
    type VARCHAR(30) NOT NULL CHECK (type IN ('internal', 'client', 'vendor', 'partner')),
    status VARCHAR(30) NOT NULL CHECK (status IN ('active', 'suspended')) DEFAULT 'active',
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_organizations_code UNIQUE (code)
);
CREATE INDEX idx_organizations_code ON organizations(code);
CREATE INDEX idx_organizations_status ON organizations(status);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    username VARCHAR(120) NOT NULL,
    email VARCHAR(180) NOT NULL,
    display_name VARCHAR(180) NOT NULL,
    type VARCHAR(30) NOT NULL CHECK (type IN ('internal', 'external', 'service')),
    status VARCHAR(30) NOT NULL CHECK (status IN ('invited', 'active', 'suspended', 'disabled')) DEFAULT 'invited',
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_users_username UNIQUE (username),
    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT fk_users_organization_id FOREIGN KEY (organization_id) REFERENCES organizations(id)
);
CREATE INDEX idx_users_organization_id ON users(organization_id);
CREATE INDEX idx_users_status ON users(status);

CREATE TABLE user_identities (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    provider VARCHAR(40) NOT NULL,
    provider_user_id VARCHAR(180) NOT NULL,
    username VARCHAR(120) NULL,
    email VARCHAR(180) NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_user_identities_provider_user UNIQUE (provider, provider_user_id),
    CONSTRAINT fk_user_identities_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_user_identities_user_id ON user_identities(user_id);
CREATE INDEX idx_user_identities_provider ON user_identities(provider);
CREATE INDEX idx_user_identities_provider_user_id ON user_identities(provider_user_id);

CREATE TABLE apps (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(120) NOT NULL,
    name VARCHAR(150) NOT NULL,
    status VARCHAR(30) NOT NULL CHECK (status IN ('active', 'inactive', 'suspended')) DEFAULT 'active',
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_apps_code UNIQUE (code)
);
CREATE INDEX idx_apps_code ON apps(code);
CREATE INDEX idx_apps_status ON apps(status);

CREATE TABLE app_clients (
    id BIGSERIAL PRIMARY KEY,
    app_id BIGINT NOT NULL,
    keycloak_client_id VARCHAR(180) NOT NULL,
    name VARCHAR(150) NOT NULL,
    status VARCHAR(30) NOT NULL CHECK (status IN ('active', 'inactive')) DEFAULT 'active',
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_app_clients_keycloak_client UNIQUE (keycloak_client_id),
    CONSTRAINT fk_app_clients_app_id FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE
);
CREATE INDEX idx_app_clients_app_id ON app_clients(app_id);
CREATE INDEX idx_app_clients_keycloak_client_id ON app_clients(keycloak_client_id);
CREATE INDEX idx_app_clients_status ON app_clients(status);

CREATE TABLE modules (
    id BIGSERIAL PRIMARY KEY,
    app_id BIGINT NOT NULL,
    code VARCHAR(120) NOT NULL,
    name VARCHAR(150) NOT NULL,
    status VARCHAR(30) NOT NULL CHECK (status IN ('active', 'inactive')) DEFAULT 'active',
    CONSTRAINT uq_modules_app_code UNIQUE (app_id, code),
    CONSTRAINT fk_modules_app_id FOREIGN KEY (app_id) REFERENCES apps(id)
);
CREATE INDEX idx_modules_app_id ON modules(app_id);
CREATE INDEX idx_modules_code ON modules(code);
CREATE INDEX idx_modules_status ON modules(status);

CREATE TABLE actions (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(120) NOT NULL,
    name VARCHAR(150) NOT NULL,
    CONSTRAINT uq_actions_code UNIQUE (code)
);
CREATE INDEX idx_actions_code ON actions(code);

CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    app_id BIGINT NOT NULL,
    module_id BIGINT NULL,
    action_id BIGINT NOT NULL,
    code VARCHAR(180) NOT NULL,
    name VARCHAR(180) NOT NULL,
    description TEXT NULL,
    risk_level VARCHAR(30) NOT NULL CHECK (risk_level IN ('low', 'medium', 'high', 'critical')) DEFAULT 'low',
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(30) NOT NULL CHECK (status IN ('active', 'inactive', 'deprecated')) DEFAULT 'active',
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_permissions_code UNIQUE (code),
    CONSTRAINT fk_permissions_app_id FOREIGN KEY (app_id) REFERENCES apps(id),
    CONSTRAINT fk_permissions_module_id FOREIGN KEY (module_id) REFERENCES modules(id),
    CONSTRAINT fk_permissions_action_id FOREIGN KEY (action_id) REFERENCES actions(id)
);
CREATE INDEX idx_permissions_app_id ON permissions(app_id);
CREATE INDEX idx_permissions_module_id ON permissions(module_id);
CREATE INDEX idx_permissions_action_id ON permissions(action_id);
CREATE INDEX idx_permissions_code ON permissions(code);
CREATE INDEX idx_permissions_status ON permissions(status);

CREATE TABLE menus (
    id BIGSERIAL PRIMARY KEY,
    app_id BIGINT NOT NULL,
    module_id BIGINT NOT NULL,
    parent_id BIGINT NULL,
    code VARCHAR(120) NOT NULL,
    name VARCHAR(150) NOT NULL,
    route_path VARCHAR(255) NOT NULL,
    sort_order BIGINT NOT NULL DEFAULT 0,
    required_permission_id BIGINT NULL,
    status VARCHAR(30) NOT NULL CHECK (status IN ('active', 'inactive')) DEFAULT 'active',
    CONSTRAINT uq_menus_app_code UNIQUE (app_id, code),
    CONSTRAINT fk_menus_app_id FOREIGN KEY (app_id) REFERENCES apps(id),
    CONSTRAINT fk_menus_module_id FOREIGN KEY (module_id) REFERENCES modules(id),
    CONSTRAINT fk_menus_parent_id FOREIGN KEY (parent_id) REFERENCES menus(id),
    CONSTRAINT fk_menus_required_permission_id FOREIGN KEY (required_permission_id) REFERENCES permissions(id)
);
CREATE INDEX idx_menus_app_id ON menus(app_id);
CREATE INDEX idx_menus_module_id ON menus(module_id);
CREATE INDEX idx_menus_parent_id ON menus(parent_id);
CREATE INDEX idx_menus_code ON menus(code);
CREATE INDEX idx_menus_required_permission_id ON menus(required_permission_id);
CREATE INDEX idx_menus_status ON menus(status);

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NULL,
    app_id BIGINT NULL,
    code VARCHAR(120) NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT NULL,
    scope VARCHAR(30) NOT NULL CHECK (scope IN ('global', 'organization', 'app')),
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(30) NOT NULL CHECK (status IN ('active', 'inactive')) DEFAULT 'active',
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_roles_scope_code UNIQUE (organization_id, app_id, code),
    CONSTRAINT fk_roles_organization_id FOREIGN KEY (organization_id) REFERENCES organizations(id),
    CONSTRAINT fk_roles_app_id FOREIGN KEY (app_id) REFERENCES apps(id)
);
CREATE INDEX idx_roles_organization_id ON roles(organization_id);
CREATE INDEX idx_roles_app_id ON roles(app_id);
CREATE INDEX idx_roles_code ON roles(code);
CREATE INDEX idx_roles_scope ON roles(scope);
CREATE INDEX idx_roles_status ON roles(status);

CREATE TABLE role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    effect VARCHAR(20) NOT NULL CHECK (effect IN ('allow', 'deny')) DEFAULT 'allow',
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_role_permissions_role_permission UNIQUE (role_id, permission_id),
    CONSTRAINT fk_role_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX idx_role_permissions_effect ON role_permissions(effect);

CREATE TABLE user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    app_id BIGINT NULL,
    organization_id BIGINT NULL,
    expires_at TIMESTAMPTZ NULL,
    assigned_by BIGINT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_user_roles_assignment UNIQUE (user_id, role_id, app_id, organization_id),
    CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_app_id FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_organization_id FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_assigned_by FOREIGN KEY (assigned_by) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_user_roles_app_id ON user_roles(app_id);
CREATE INDEX idx_user_roles_organization_id ON user_roles(organization_id);

CREATE TABLE teams (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    code VARCHAR(120) NOT NULL,
    name VARCHAR(150) NOT NULL,
    status VARCHAR(30) NOT NULL CHECK (status IN ('active', 'inactive')) DEFAULT 'active',
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_teams_org_code UNIQUE (organization_id, code),
    CONSTRAINT fk_teams_organization_id FOREIGN KEY (organization_id) REFERENCES organizations(id)
);
CREATE INDEX idx_teams_organization_id ON teams(organization_id);
CREATE INDEX idx_teams_code ON teams(code);
CREATE INDEX idx_teams_status ON teams(status);

CREATE TABLE team_members (
    id BIGSERIAL PRIMARY KEY,
    team_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role_in_team VARCHAR(80) NULL,
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_team_members_team_user UNIQUE (team_id, user_id),
    CONSTRAINT fk_team_members_team_id FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    CONSTRAINT fk_team_members_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_team_members_team_id ON team_members(team_id);
CREATE INDEX idx_team_members_user_id ON team_members(user_id);

CREATE TABLE team_roles (
    id BIGSERIAL PRIMARY KEY,
    team_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    app_id BIGINT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_team_roles_team_role_app UNIQUE (team_id, role_id, app_id),
    CONSTRAINT fk_team_roles_team_id FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    CONSTRAINT fk_team_roles_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_team_roles_app_id FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE
);
CREATE INDEX idx_team_roles_team_id ON team_roles(team_id);
CREATE INDEX idx_team_roles_role_id ON team_roles(role_id);
CREATE INDEX idx_team_roles_app_id ON team_roles(app_id);

CREATE TABLE user_permission_overrides (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    effect VARCHAR(20) NOT NULL CHECK (effect IN ('allow', 'deny')),
    reason TEXT NOT NULL,
    expires_at TIMESTAMPTZ NULL,
    created_by BIGINT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_user_permission_overrides_user_permission UNIQUE (user_id, permission_id),
    CONSTRAINT fk_user_permission_overrides_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_permission_overrides_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_permission_overrides_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_user_permission_overrides_user_id ON user_permission_overrides(user_id);
CREATE INDEX idx_user_permission_overrides_permission_id ON user_permission_overrides(permission_id);
CREATE INDEX idx_user_permission_overrides_effect ON user_permission_overrides(effect);

CREATE TABLE access_cache_versions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    app_id BIGINT NOT NULL,
    version BIGINT NOT NULL DEFAULT 1,
    updated_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_access_cache_versions_user_app UNIQUE (user_id, app_id),
    CONSTRAINT fk_access_cache_versions_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_access_cache_versions_app_id FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE
);
CREATE INDEX idx_access_cache_versions_user_id ON access_cache_versions(user_id);
CREATE INDEX idx_access_cache_versions_app_id ON access_cache_versions(app_id);

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NULL,
    app_id BIGINT NULL,
    actor_user_id BIGINT NULL,
    target_user_id BIGINT NULL,
    permission_id BIGINT NULL,
    action VARCHAR(120) NOT NULL,
    resource_type VARCHAR(120) NOT NULL,
    resource_id VARCHAR(120) NULL,
    result VARCHAR(30) NOT NULL CHECK (result IN ('allowed', 'denied', 'changed')),
    metadata_json JSONB NULL DEFAULT '{}'::jsonb,
    ip_address VARCHAR(80) NULL,
    user_agent TEXT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT fk_audit_logs_organization_id FOREIGN KEY (organization_id) REFERENCES organizations(id),
    CONSTRAINT fk_audit_logs_app_id FOREIGN KEY (app_id) REFERENCES apps(id),
    CONSTRAINT fk_audit_logs_actor_user_id FOREIGN KEY (actor_user_id) REFERENCES users(id),
    CONSTRAINT fk_audit_logs_target_user_id FOREIGN KEY (target_user_id) REFERENCES users(id),
    CONSTRAINT fk_audit_logs_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id)
);
CREATE INDEX idx_audit_logs_organization_id ON audit_logs(organization_id);
CREATE INDEX idx_audit_logs_app_id ON audit_logs(app_id);
CREATE INDEX idx_audit_logs_actor_user_id ON audit_logs(actor_user_id);
CREATE INDEX idx_audit_logs_target_user_id ON audit_logs(target_user_id);
CREATE INDEX idx_audit_logs_permission_id ON audit_logs(permission_id);
CREATE INDEX idx_audit_logs_resource_id ON audit_logs(resource_id);
CREATE INDEX idx_audit_logs_result ON audit_logs(result);

CREATE TABLE app_permission_manifests (
    id BIGSERIAL PRIMARY KEY,
    app_id BIGINT NOT NULL,
    version BIGINT NOT NULL,
    checksum VARCHAR(128) NOT NULL,
    manifest_json JSONB NULL,
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    activated_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_app_permission_manifests_app_version UNIQUE (app_id, version),
    CONSTRAINT fk_app_permission_manifests_app_id FOREIGN KEY (app_id) REFERENCES apps(id)
);
CREATE INDEX idx_app_permission_manifests_app_id ON app_permission_manifests(app_id);

-- +goose Down
DROP TABLE IF EXISTS app_permission_manifests;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS access_cache_versions;
DROP TABLE IF EXISTS user_permission_overrides;
DROP TABLE IF EXISTS team_roles;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS menus;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS actions;
DROP TABLE IF EXISTS modules;
DROP TABLE IF EXISTS app_clients;
DROP TABLE IF EXISTS apps;
DROP TABLE IF EXISTS user_identities;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS organizations;
