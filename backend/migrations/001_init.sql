CREATE TABLE IF NOT EXISTS sys_depts (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  created_by BIGINT NOT NULL DEFAULT 0,
  updated_by BIGINT NOT NULL DEFAULT 0,
  parent_id BIGINT NOT NULL DEFAULT 0,
  name VARCHAR(64) NOT NULL,
  sort INTEGER NOT NULL DEFAULT 0,
  status VARCHAR(16) NOT NULL DEFAULT 'enabled'
);

CREATE TABLE IF NOT EXISTS sys_users (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  created_by BIGINT NOT NULL DEFAULT 0,
  updated_by BIGINT NOT NULL DEFAULT 0,
  username VARCHAR(64) NOT NULL UNIQUE,
  nickname VARCHAR(64) NOT NULL,
  password VARCHAR(255) NOT NULL,
  email VARCHAR(128),
  mobile VARCHAR(32),
  avatar_id BIGINT,
  status VARCHAR(16) NOT NULL DEFAULT 'enabled',
  dept_id BIGINT
);

CREATE TABLE IF NOT EXISTS sys_roles (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  created_by BIGINT NOT NULL DEFAULT 0,
  updated_by BIGINT NOT NULL DEFAULT 0,
  code VARCHAR(64) NOT NULL UNIQUE,
  name VARCHAR(64) NOT NULL,
  sort INTEGER NOT NULL DEFAULT 0,
  status VARCHAR(16) NOT NULL DEFAULT 'enabled',
  data_scope VARCHAR(32) NOT NULL DEFAULT 'self'
);

CREATE TABLE IF NOT EXISTS sys_menus (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  created_by BIGINT NOT NULL DEFAULT 0,
  updated_by BIGINT NOT NULL DEFAULT 0,
  parent_id BIGINT NOT NULL DEFAULT 0,
  name VARCHAR(64) NOT NULL,
  title VARCHAR(64) NOT NULL,
  type VARCHAR(16) NOT NULL,
  path VARCHAR(255),
  component VARCHAR(255),
  icon VARCHAR(64),
  permission VARCHAR(128),
  hidden BOOLEAN NOT NULL DEFAULT FALSE,
  sort INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS sys_user_roles (
  user_id BIGINT NOT NULL REFERENCES sys_users(id) ON DELETE CASCADE,
  role_id BIGINT NOT NULL REFERENCES sys_roles(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS sys_role_menus (
  role_id BIGINT NOT NULL REFERENCES sys_roles(id) ON DELETE CASCADE,
  menu_id BIGINT NOT NULL REFERENCES sys_menus(id) ON DELETE CASCADE,
  PRIMARY KEY (role_id, menu_id)
);

CREATE TABLE IF NOT EXISTS sys_files (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  original_name VARCHAR(255) NOT NULL,
  bucket VARCHAR(128) NOT NULL,
  object_key VARCHAR(512) NOT NULL,
  content_type VARCHAR(128),
  size BIGINT NOT NULL DEFAULT 0,
  created_by BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS sys_login_logs (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  username VARCHAR(64),
  ip VARCHAR(64),
  status VARCHAR(32),
  message VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS sys_operation_logs (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id BIGINT NOT NULL DEFAULT 0,
  module VARCHAR(64),
  action VARCHAR(64),
  method VARCHAR(16),
  path VARCHAR(255),
  ip VARCHAR(64),
  status INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_sys_users_deleted_at ON sys_users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_sys_roles_deleted_at ON sys_roles(deleted_at);
CREATE INDEX IF NOT EXISTS idx_sys_menus_deleted_at ON sys_menus(deleted_at);
CREATE INDEX IF NOT EXISTS idx_sys_menus_permission ON sys_menus(permission);

INSERT INTO sys_depts(id, parent_id, name, sort, status)
VALUES (1, 0, '总部', 1, 'enabled')
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_users(id, username, nickname, password, email, status, dept_id)
VALUES (1, 'admin', '系统管理员', '$2b$12$izSwz7dsGtep4rwOcn3II.dT5pjqynrCHwUysW8fu0IgnnHafo5Wq', 'admin@example.com', 'enabled', 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_roles(id, code, name, sort, status, data_scope)
VALUES (1, 'admin', '超级管理员', 1, 'enabled', 'all')
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_user_roles(user_id, role_id)
VALUES (1, 1)
ON CONFLICT DO NOTHING;

INSERT INTO sys_menus(id, parent_id, name, title, type, path, component, icon, permission, sort) VALUES
(1, 0, 'Dashboard', '仪表盘', 'menu', '/dashboard', 'dashboard/index', 'DashboardOutlined', '', 1),
(2, 0, 'System', '系统管理', 'catalog', '/system', '', 'SettingOutlined', '', 10),
(3, 2, 'User', '用户管理', 'menu', '/system/user', 'system/user/index', 'UserOutlined', 'system:user:list', 1),
(4, 2, 'Role', '角色管理', 'menu', '/system/role', 'system/role/index', 'TeamOutlined', 'system:role:list', 2),
(5, 2, 'Menu', '菜单管理', 'menu', '/system/menu', 'system/menu/index', 'MenuOutlined', 'system:menu:list', 3),
(6, 2, 'Dept', '部门管理', 'menu', '/system/dept', 'system/dept/index', 'ApartmentOutlined', 'system:dept:list', 4),
(7, 0, 'File', '文件管理', 'menu', '/file', 'file/index', 'FileOutlined', 'file:read', 20),
(8, 0, 'Audit', '审计日志', 'catalog', '/audit', '', 'AuditOutlined', '', 30),
(9, 8, 'LoginLog', '登录日志', 'menu', '/audit/login-log', 'audit/login-log/index', 'LoginOutlined', 'audit:login-log:list', 1),
(10, 8, 'OperationLog', '操作日志', 'menu', '/audit/operation-log', 'audit/operation-log/index', 'ProfileOutlined', 'audit:operation-log:list', 2),
(101, 3, 'UserCreate', '新增用户', 'button', '', '', '', 'system:user:create', 1),
(102, 3, 'UserUpdate', '编辑用户', 'button', '', '', '', 'system:user:update', 2),
(103, 3, 'UserDelete', '删除用户', 'button', '', '', '', 'system:user:delete', 3),
(111, 4, 'RoleCreate', '新增角色', 'button', '', '', '', 'system:role:create', 1),
(112, 4, 'RoleUpdate', '编辑角色', 'button', '', '', '', 'system:role:update', 2),
(113, 4, 'RoleDelete', '删除角色', 'button', '', '', '', 'system:role:delete', 3),
(121, 7, 'FileUpload', '上传文件', 'button', '', '', '', 'file:upload', 1),
(122, 7, 'FileDelete', '删除文件', 'button', '', '', '', 'file:delete', 2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_role_menus(role_id, menu_id)
SELECT 1, id FROM sys_menus
ON CONFLICT DO NOTHING;

SELECT setval('sys_depts_id_seq', (SELECT MAX(id) FROM sys_depts));
SELECT setval('sys_users_id_seq', (SELECT MAX(id) FROM sys_users));
SELECT setval('sys_roles_id_seq', (SELECT MAX(id) FROM sys_roles));
SELECT setval('sys_menus_id_seq', (SELECT MAX(id) FROM sys_menus));
