INSERT INTO sys_menus(id, parent_id, name, title, type, path, component, icon, permission, sort)
VALUES (122, 7, 'FileDelete', '删除文件', 'button', '', '', '', 'file:delete', 2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_role_menus(role_id, menu_id)
SELECT 1, 122
WHERE EXISTS (SELECT 1 FROM sys_roles WHERE id = 1)
ON CONFLICT DO NOTHING;

SELECT setval('sys_menus_id_seq', (SELECT MAX(id) FROM sys_menus));
